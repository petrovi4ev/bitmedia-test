package apiserver

import (
	"context"
	"encoding/json"
	mdbp "github.com/gobeam/mongo-go-pagination"
	"github.com/gorilla/mux"
	"github.com/petrovi4ev/bitmedia-test/internal/config"
	"github.com/petrovi4ev/bitmedia-test/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strconv"
	"time"
)

type APIServer struct {
	DB         *mongo.Client
	Config     config.Config
	Collection *mongo.Collection
}

type Pagination struct {
	CurrentPage int
	PageCount   int
	TotalCount  int
	PerPage     int
}

func New(db *mongo.Client, cfg config.Config) *APIServer {
	return &APIServer{
		DB:         db,
		Config:     cfg,
		Collection: db.Database(cfg.DbName).Collection("users"),
	}
}

func (s *APIServer) NewPagination() *Pagination {
	perPage, err := strconv.Atoi(s.Config.PaginationPerPage)
	check(err)
	return &Pagination{
		CurrentPage: 1,
		PageCount:   0,
		TotalCount:  0,
		PerPage:     perPage,
	}
}

func (s *APIServer) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/users", s.AllUsersEndPoint).Methods("GET")
	r.HandleFunc("/users", s.CreateUserEndPoint).Methods("POST")
	r.HandleFunc("/users/{id}", s.UpdateUserEndPoint).Methods("PUT")
	r.HandleFunc("/users/{id}", s.DeleteUserEndPoint).Methods("DELETE")
	r.HandleFunc("/users/{id}", s.FindUserEndpoint).Methods("GET")
	if err := http.ListenAndServe(":"+s.Config.ServerPort, r); err != nil {
		log.Fatal(err)
	}
}

func (s *APIServer) AllUsersEndPoint(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	perPageParam := r.URL.Query().Get("per-page")
	currentPage := 1
	if len(page) > 0 {
		res, err := strconv.Atoi(page)
		check(err)
		currentPage = res
	}
	perPage, err := strconv.Atoi(s.Config.PaginationPerPage)
	if len(perPageParam) > 0 {
		perPage, err = strconv.Atoi(perPageParam)
		check(err)
	}
	paginatedData, err := mdbp.New(s.Collection).Limit(int64(perPage)).Page(int64(currentPage)).Filter(bson.D{}).Find()
	check(err)

	var lists []model.User
	if paginatedData != nil {
		for _, raw := range paginatedData.Data {
			var product *model.User
			if marshallErr := bson.Unmarshal(raw, &product); marshallErr == nil {
				lists = append(lists, *product)
			}
		}
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := &Pagination{
		CurrentPage: currentPage,
		PageCount:   int(paginatedData.Pagination.TotalPage),
		TotalCount:  int(paginatedData.Pagination.Total),
		PerPage:     perPage,
	}
	response := map[string]interface{}{
		"users":      lists,
		"pagination": *pagination,
	}
	pagination = s.NewPagination()
	respondWithJson(w, http.StatusOK, response)
}

func (s *APIServer) FindUserEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	objID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Println(err)
		return
	}
	var user model.User
	filter := bson.M{"_id": objID}
	err = s.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	respondWithJson(w, http.StatusOK, user)
}

func (s *APIServer) CreateUserEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user.ID = primitive.NewObjectID()
	if s.CheckEmail(user.Email) {
		respondWithError(w, http.StatusBadRequest, "Email already exist.")
		return
	}
	result, err := s.Collection.InsertOne(ctx, user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	objectID := result.InsertedID.(primitive.ObjectID)
	log.Println(objectID)
	respondWithJson(w, http.StatusCreated, user)
}

func (s *APIServer) UpdateUserEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	objID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Println(err)
		return
	}
	var user model.User
	filter := bson.M{"_id": objID}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var userMail model.User
	filterUnique := bson.M{
		"email": user.Email,
	}
	err = s.Collection.FindOne(ctx, filterUnique).Decode(&userMail)
	check(err)
	log.Println(userMail.ID.String())
	if !userMail.ID.IsZero() {
		if userMail.ID != objID {
			respondWithError(w, http.StatusBadRequest, "Email already exist.")
			return
		}
	}

	update := bson.M{
		"$set": bson.M{
			"last_name":  user.LastName,
			"gender":     user.Gender,
			"email":      user.Email,
			"country":    user.Country,
			"city":       user.City,
			"birth_date": user.BirthDate,
		},
	}
	resultUpdate, err := s.Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println(resultUpdate.ModifiedCount)
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func (s *APIServer) DeleteUserEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	objID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Println(err)
		return
	}
	resultDelete, err := s.Collection.DeleteOne(ctx, bson.M{"_id": objID})
	if resultDelete != nil {
		log.Println(resultDelete.DeletedCount)
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (s *APIServer) CheckEmail(e string) bool {
	filter := bson.M{"email": e}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	c, err := s.Collection.CountDocuments(ctx, filter)
	check(err)
	if c > 0 {
		return true
	}

	return false
}
