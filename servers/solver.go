package servers

import (
	"github.com/go-chi/chi/v5"
	"github.com/thinkingatoms/apibase/ez"
	"github.com/thinkingatoms/apibase/models"
	"github.com/thinkingatoms/apibase/servers"
	_models "github.com/thinkingatoms/wordgames/models"
	errors "golang.org/x/xerrors"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

type solverService struct {
	models.JWTIssuer
	server     *servers.Server
	wordGames  *_models.WordGames
	cache      *models.TenureCache
	adminRoles map[string]bool
	newWords   []string
}

func RegisterWordGamesSolver(server *servers.Server) error {
	cache := server.GetCache()
	if cache == nil {
		return errors.New("no cache")
	}
	name := "games"
	if !server.HasSubConfig(name) {
		return nil
	}
	s := solverService{
		JWTIssuer:  models.NewJWTIssuer(server.GetSecret),
		server:     server,
		cache:      cache,
		adminRoles: map[string]bool{"admin": true},
		wordGames:  _models.WordGamesFromConfig(server.GetSubConfig(name), cache),
	}
	s.EnrichRouter(server.Router)
	return nil
}

func (self *solverService) EnrichRouter(router *chi.Mux) {
	router.Get("/", ez.StaticMsgHandler(strconv.Itoa(len(self.wordGames.WordBank.Infos))))
	router.Get("/killerwasp", self.killerWaspHandler)
}

var Slash = ([]rune("/"))[0]

func (self *solverService) killerWaspHandler(w http.ResponseWriter, r *http.Request) {
	s := strings.ToUpper(r.URL.Query().Get("s"))
	req := make([]rune, 0)
	opt := make([]rune, 0)
	last := rune(0)
	for i, w := 0, 0; i < len(s); i += w {
		r, width := utf8.DecodeRuneInString(s[i:])
		if r == Slash {
		} else if last == Slash {
			req = append(req, r)
		} else {
			opt = append(opt, r)
		}
		last = r
		w = width
	}
	words := _models.NewKillerWasp(req, opt, 0).Solve(self.wordGames.WordBank)
	ez.WriteObjectAsJSON(w, r, words)
}
