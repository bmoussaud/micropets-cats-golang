package cats

import (
	"fmt"

	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	. "moussaud.org/cats/internal"
)

// Cat Struct
type Cat struct {
	Index int
	Name  string
	Kind  string
	Age   int
	URL   string
	From  string
	URI   string
}

// Cats Struct
type Cats struct {
	Total    int
	Hostname string
	Cats     []Cat `json:"Pets"`
}

var calls = 0

var shift = 0

func setupResponse(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

}

func db() Cats {
	cat1 := Cat{20, "Orphée", "Persan", 12, "https://cdn.pixabay.com/photo/2020/02/29/13/51/cat-4890133_960_720.jpg", GlobalConfig.Service.From, "/cats/v1/data/0"}
	cat2 := Cat{21, "Pirouette", "Bengal", 1, "https://upload.wikimedia.org/wikipedia/commons/thumb/b/ba/Paintedcats_Red_Star_standing.jpg/934px-Paintedcats_Red_Star_standing.jpg", GlobalConfig.Service.From, "/cats/v1/data/1"}
	cat3 := Cat{22, "Pamina", "Angora", 120, "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a5/Turkish_Angora_Odd-Eyed.jpg/440px-Turkish_Angora_Odd-Eyed.jpg", GlobalConfig.Service.From, "/cats/v1/data/2"}
	cat4 := Cat{23, "Tommy Lee", "Siamois", 120, "https://www.woopets.fr/assets/races/home/siamois-124x153.jpg", GlobalConfig.Service.From, "/cats/v1/data/3"}
	cats := Cats{4, "Unknown", []Cat{cat1, cat2, cat3, cat4}}
	host, err := os.Hostname()

	if err != nil {
		cats.Hostname = "Unknown"
	} else {
		cats.Hostname = host
	}

	return cats
}

func db_authentication(r *http.Request) {
	//span := NewServerSpan(r, "db_authentication")
	//defer span.Finish()

	NewTrace(r.Context(), "db_authentication")

	if GlobalConfig.Service.Delay.Period > 0 {
		y := float64(calls+shift) * math.Pi / float64(2*GlobalConfig.Service.Delay.Period)
		sin_y := math.Sin(y)
		abs_y := math.Abs(sin_y)
		sleep := int(abs_y * GlobalConfig.Service.Delay.Amplitude * 1000.0)
		fmt.Printf("waitTime %d - %f - %f - %f  -> sleep %d ms  \n", calls, y, math.Sin(y), abs_y, sleep)
		start := time.Now()
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		elapsed := time.Since(start)
		fmt.Printf("Current Unix Time: %s\n", elapsed)
	}
}

func single(c *gin.Context) {

	setupResponse(c)
	time.Sleep(time.Duration(10) * time.Millisecond)

	db_authentication(c.Request)

	cats := db()

	strId := c.Param("id")
	id, err := strconv.Atoi(strId)

	if err != nil {
		fmt.Println("Error during conversion")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error during conversion %s", strId)})
		return
	}

	fmt.Printf("ID %d", id)
	if id >= len(cats.Cats) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("invalid index %d", id)})
	} else {
		element := cats.Cats[id]
		element.From = GlobalConfig.Service.From
		fmt.Println(element)
		c.IndentedJSON(http.StatusOK, element)
	}
}

func index(c *gin.Context) {

	setupResponse(c)
	time.Sleep(time.Duration(10) * time.Millisecond)

	db_authentication(c.Request)

	cats := db()

	for i := 1; i < cats.Total; i++ {
		cats.Cats[i].From = GlobalConfig.Service.From
	}

	calls = calls + 1
	if GlobalConfig.Service.Mode == "RANDOM_NUMBER" {
		total := rand.Intn(cats.Total) + 1
		//fmt.Printf("reduce results to total %d/%d\n", total, cats.Total)
		for i := 1; i < total; i++ {
			cats.Cats = cats.Cats[:len(cats.Cats)-1]
			cats.Total = cats.Total - 1
		}
	}

	if GlobalConfig.Service.FrequencyError > 0 && calls%GlobalConfig.Service.FrequencyError == 0 {
		fmt.Printf("Fails this call (%d)", calls)
		//otrext.Error.Set(span, true)
		//span.LogFields(otrlog.String("error.kind", "Unexpected Error when querying the cats repository"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected Error when querying the cats repository"})

	} else {
		c.IndentedJSON(http.StatusOK, cats)
	}
}

// GetLocation returns the full path of the config file based on the current executable location or using SERVICE_CONFIG_DIR env
func GetLocation(file string) string {
	if serviceConfigDirectory := os.Getenv("SERVICE_CONFIG_DIR"); serviceConfigDirectory != "" {
		fmt.Printf("Load configuration from %s\n", serviceConfigDirectory)
		return filepath.Join(serviceConfigDirectory, file)
	} else {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		return filepath.Join(exPath, file)
	}
}

func readiness_and_liveness(c *gin.Context) {
	//span := NewServerSpan(r, "readiness_and_liveness")
	//defer span.Finish()
	//defer span.Finish()
	//NewTrace(r.Context(), "readiness_and_liveness")
	c.String(http.StatusOK, "OK\n")
}

func Start() {
	config := LoadConfiguration()

	r := gin.Default()
	r.Use(otelgin.Middleware("otel-otlp-go-service"))

	r.GET("/cats/v1/data", index)
	r.GET("/cats/v1/data/:id", single)

	r.GET("/cats/liveness", readiness_and_liveness)
	r.GET("/cats/readiness", readiness_and_liveness)

	r.GET("/liveness", readiness_and_liveness)
	r.GET("/readiness", readiness_and_liveness)

	r.GET("/", index)

	rand.Seed(time.Now().UnixNano())
	shift = rand.Intn(100)

	fmt.Printf("******* Starting to the cats service on port %s, mode %s\n", config.Service.Port, config.Service.Mode)
	fmt.Printf("******* Delay Period %d Amplitude %f shift %d \n", config.Service.Delay.Period, config.Service.Delay.Amplitude, shift)
	fmt.Printf("******* Frequency Error %d\n", config.Service.FrequencyError)

	r.Run(config.Service.Port)
}
