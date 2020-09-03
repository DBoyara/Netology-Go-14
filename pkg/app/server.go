package app

import (
	"encoding/json"
	"github.com/DBoyara/Netology-Go-14/pkg/app/dto"
	"github.com/DBoyara/Netology-Go-14/pkg/card"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	cardSvc *card.Service
	mux     *http.ServeMux
}

func NewServer(cardSvc *card.Service, mux *http.ServeMux) *Server {
	return &Server{cardSvc: cardSvc, mux: mux}
}

func (s *Server) Init() {
	s.mux.HandleFunc("/getCards", s.getCards)
	s.mux.HandleFunc("/addCard", s.addCard)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) getCards(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		dtos := dto.CardErrDTO{Err: card.ErrNotSpecifiedUserId.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	intUserId, err := strconv.Atoi(userId)
	if err != nil {
		dtos := dto.CardErrDTO{Err: card.ErrUserDoesNotExist.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	cards, _ := s.cardSvc.All(card.UserID(intUserId))
	var dtos []*dto.CardDTO

	for _, c := range cards {
		dtos = append(
			dtos,
			&dto.CardDTO{
				Id:     c.Id,
				Number: c.Number,
				Type:   c.Type,
				Issuer: c.Issuer,
				UserId: card.UserID(intUserId),
			})
	}
	jsonResponse(w, r, dtos)
}

func (s *Server) addCard(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		dtos := dto.CardErrDTO{Err: err.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	params := &dto.CardDTO{}
	err = json.Unmarshal(body, params)
	if err != nil {
		dtos := dto.CardErrDTO{Err: err.Error()}
		jsonResponse(w, r, dtos)
		return
	}
	newCard, err := s.cardSvc.Add(params.UserId, params.Type, params.Issuer)

	if err != nil {
		dtos := dto.CardErrDTO{Err: err.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	var dtos []*dto.CardDTO
	dtos = append(dtos,
		&dto.CardDTO{
			Id:     newCard.Id,
			UserId: params.UserId,
			Number: newCard.Number,
			Type:   newCard.Type,
			Issuer: newCard.Issuer,
		})

	jsonResponse(w, r, dtos)
}

func jsonResponse(w http.ResponseWriter, r *http.Request, dtos... interface{}) {
	respBody, err := json.Marshal(dtos)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(respBody)
	if err != nil {
		log.Println(err)
		return
	}
}
