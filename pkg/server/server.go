package app

import (
	"context"
	"encoding/json"
	"github.com/DBoyara/Netology-Go-14/pkg/app/dto"
	"github.com/DBoyara/Netology-Go-14/pkg/card"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	cardSvc *card.Service
	mux     *http.ServeMux
	ctx     context.Context
	conn    *pgxpool.Conn
}

func NewServer(cardSvc *card.Service, mux *http.ServeMux, ctx context.Context, conn *pgxpool.Conn) *Server {
	return &Server{cardSvc: cardSvc, mux: mux, ctx: ctx, conn: conn}
}

func (s *Server) Init() {
	s.mux.HandleFunc("/getCards", s.getCards)
	s.mux.HandleFunc("/addCard", s.addCard)
	s.mux.HandleFunc("/getTransactions", s.getTransactions)
	s.mux.HandleFunc("/getMostPaid", s.getMostPaid)
	s.mux.HandleFunc("/getMostOftenBought", s.getMostOftenBought)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) getCards(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		dtos := dto.ErrDTO{Err: card.ErrNotSpecifiedUserId.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	intUserId, err := strconv.Atoi(userId)
	if err != nil {
		dtos := dto.ErrDTO{Err: card.ErrUserDoesNotExist.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	cards, _ := s.cardSvc.All(card.UserID(intUserId))
	var dtos []*dto.CardDTO

	for _, c := range cards {
		dtos = append(
			dtos,
			&dto.CardDTO{
				Id:      c.Id,
				Number:  c.Number,
				Type:    c.Type,
				Issuer:  c.Issuer,
				OwnerId: card.UserID(intUserId),
				Balance: c.Balance,
				Status:  c.Status,
				Created: c.Created,
			})
	}
	jsonResponse(w, r, dtos)
}

func (s *Server) addCard(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		dtos := dto.ErrDTO{Err: err.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	params := &dto.CardDTO{}
	err = json.Unmarshal(body, params)
	if err != nil {
		dtos := dto.ErrDTO{Err: err.Error()}
		jsonResponse(w, r, dtos)
		return
	}
	newCard, err := s.cardSvc.Add(params.OwnerId, params.Type, params.Issuer)

	if err != nil {
		dtos := dto.ErrDTO{Err: err.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	var dtos []*dto.CardDTO
	dtos = append(dtos,
		&dto.CardDTO{
			Id:      newCard.Id,
			OwnerId: params.OwnerId,
			Number:  newCard.Number,
			Type:    newCard.Type,
			Issuer:  newCard.Issuer,
			Balance: newCard.Balance,
			Status:  newCard.Status,
			Created: newCard.Created,
		})

	jsonResponse(w, r, dtos)
}

func jsonResponse(w http.ResponseWriter, r *http.Request, dtos ...interface{}) {
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

func (s *Server) getTransactions(w http.ResponseWriter, r *http.Request) {
	cardId := r.URL.Query().Get("cardId")
	if cardId == "" {
		dtos := dto.ErrDTO{Err: card.ErrCardNotFound.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	intCardId, err := strconv.Atoi(cardId)
	var transactions []*dto.TransactionDTO
	rows, err := s.conn.Query(s.ctx, `
		SELECT id, card_id, amount, created, status, mcc_id, description, icon_id
		FROM transactions 
		WHERE card_id = $1
		LIMIT 5
	`, intCardId)

	defer rows.Close()

	for rows.Next() {
		tr := &dto.TransactionDTO{}
		err = rows.Scan(&tr.Id, &tr.CardId, &tr.Amount, &tr.Created, &tr.Status, &tr.MccId, &tr.Description, &tr.IconId)
		if err != nil {
			dtos := dto.ErrDTO{Err: card.ErrNoRowExec.Error()}
			jsonResponse(w, r, dtos)
			return
		}
		transactions = append(transactions, tr)
	}

	jsonResponse(w, r, transactions)
}

func (s *Server) getMostPaid(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		dtos := dto.ErrDTO{Err: card.ErrNotSpecifiedUserId.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	intUserId, err := strconv.Atoi(userId)
	if err != nil {
		dtos := dto.ErrDTO{Err: card.ErrUserDoesNotExist.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	mostPaid := &dto.MostPaid{}
	err = s.conn.QueryRow(s.ctx, `
		SELECT m.id, SUM(t.amount) AS sum_transactions, m.description
		FROM cards c
		JOIN transactions t ON c.id = t.card_id
		JOIN mcc m ON t.mcc_id = m.id
		WHERE c.owner_id = $1 AND t.amount < 0
		GROUP BY m.id, m.description
		ORDER BY sum_transactions
		LIMIT 1;
	`, intUserId).Scan(&mostPaid.MccId, &mostPaid.Amount, &mostPaid.Description)
	if err != nil {
		if err != pgx.ErrNoRows {
			dtos := dto.ErrDTO{Err: card.ErrNoRowExec.Error()}
			jsonResponse(w, r, dtos)
			return
		}
	}

	jsonResponse(w, r, mostPaid)
}

func (s *Server) getMostOftenBought(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		dtos := dto.ErrDTO{Err: card.ErrNotSpecifiedUserId.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	intUserId, err := strconv.Atoi(userId)
	if err != nil {
		dtos := dto.ErrDTO{Err: card.ErrUserDoesNotExist.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	mostOftenBought := &dto.MostOftenBought{}
	err = s.conn.QueryRow(s.ctx, `
		SELECT t.mcc_id, count(*) AS "count", m.description
		FROM cards c
		JOIN transactions t ON c.id = t.card_id
		JOIN mcc m ON t.mcc_id = m.id
		WHERE c.owner_id = $1 AND t.amount < 0
		GROUP BY t.mcc_id, m.description
		ORDER BY "count" DESC
		LIMIT 1;
	`, intUserId).Scan(&mostOftenBought.MccId, &mostOftenBought.Count, &mostOftenBought.Description)
	if err != nil {
		if err != pgx.ErrNoRows {
			dtos := dto.ErrDTO{Err: pgx.ErrNoRows.Error()}
			jsonResponse(w, r, dtos)
			return
		}
	}

	jsonResponse(w, r, mostOftenBought)
}
