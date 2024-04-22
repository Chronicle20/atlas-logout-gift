package gift

import (
	"atlas-lgs/rest"
	"github.com/gorilla/mux"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

const (
	getWorldGifts = "get_world_gifts"
)

func InitResource(si jsonapi.ServerInformation) func(router *mux.Router, l logrus.FieldLogger, db *gorm.DB) {
	return func(router *mux.Router, l logrus.FieldLogger, db *gorm.DB) {
		dRouter := router.PathPrefix("/worlds/{worldId}/gifts").Subrouter()
		dRouter.HandleFunc("", registerGetWorldGifts(si)(l, db)).Methods(http.MethodGet)
	}
}

type WorldIdHandler func(worldId int32) http.HandlerFunc

func WorldParseId(l logrus.FieldLogger, next WorldIdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		worldId, err := strconv.Atoi(mux.Vars(r)["worldId"])
		if err != nil {
			l.Errorf("Unable to properly parse worldId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(int32(worldId))(w, r)
	}
}

func registerGetWorldGifts(si jsonapi.ServerInformation) func(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return func(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
		return rest.RetrieveSpan(getWorldGifts, func(span opentracing.Span) http.HandlerFunc {
			return WorldParseId(l, func(worldId int32) http.HandlerFunc {
				return handleGetWorldGifts(si)(l, db)(span)(worldId)
			})
		})
	}
}

func handleGetWorldGifts(si jsonapi.ServerInformation) func(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(worldId int32) http.HandlerFunc {
	return func(l logrus.FieldLogger, db *gorm.DB) func(span opentracing.Span) func(worldId int32) http.HandlerFunc {
		return func(span opentracing.Span) func(worldId int32) http.HandlerFunc {
			return func(worldId int32) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					gs, err := GetForWorld(l, db)(worldId)
					if err != nil {
						l.WithError(err).Errorf("Unable to locate gifts for world [%d].", worldId)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					res, err := jsonapi.MarshalWithURLs(TransformAll(gs), si)
					if err != nil {
						l.WithError(err).Errorf("Unable to marshal models.")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					_, err = w.Write(res)
					if err != nil {
						l.WithError(err).Errorf("Unable to write response.")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
			}
		}
	}
}
