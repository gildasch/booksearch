package booksearch

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoogleVision(t *testing.T) {
	b, err := ioutil.ReadFile("test_image.jpg")
	require.NoError(t, err)

	name, short, description, err := GoogleVision(b)
	require.NoError(t, err)

	assert.Equal(t, "Logicomix", name)
	assert.Equal(t, "Roman d'Apóstolos Doxiádis et Christos Papadimitriou", short)
	assert.Equal(t, "Logicomix est un roman graphique sur la quête des fondements des mathématiques. Il est scénarisé par l'écrivain Apóstolos K. Doxiàdis et le professeur et chercheur en informatique théorique Christos H. Papadimitriou. ", description)
}

func TestEntity(t *testing.T) {
	name, short, description, err := entity("/m/05q8g20")
	require.NoError(t, err)

	assert.Equal(t, "Logicomix", name)
	assert.Equal(t, "Roman d'Apóstolos Doxiádis et Christos Papadimitriou", short)
	assert.Equal(t, "Logicomix est un roman graphique sur la quête des fondements des mathématiques. Il est scénarisé par l'écrivain Apóstolos K. Doxiàdis et le professeur et chercheur en informatique théorique Christos H. Papadimitriou. ", description)
}

const expectedResponse = `{
  "responses": [
    {
      "webDetection": {
        "webEntities": [
          {
            "entityId": "/m/057_2_",
            "score": 10.330501,
            "description": "Apostolos Doxiadis"
          },
          {
            "entityId": "/m/05q8g20",
            "score": 1.0599,
            "description": "Logicomix"
          },
          {
            "entityId": "/m/0bt_c3",
            "score": 0.5113,
            "description": "Book"
          },
          {
            "entityId": "/m/03gq5hm",
            "score": 0.500747,
            "description": "Font"
          }
        ],
        "visuallySimilarImages": [
          {
            "url": "https://anniesbookz.files.wordpress.com/2014/10/img_1945.jpg"
          },
          {
            "url": "https://cloud10.todocoleccion.online/comics/tc/2016/04/05/18/55884219.jpg"
          },
          {
            "url": "https://lefilrouge3.files.wordpress.com/2014/12/logicomix.jpg"
          },
          {
            "url": "https://http2.mlstatic.com/logicomix-riverside-agency-salamandra-graphic-D_NQ_NP_961233-MLA25749733682_072017-F.jpg"
          },
          {
            "url": "https://farm8.staticflickr.com/7224/7045642755_91fcd3e451_b.jpg"
          },
          {
            "url": "http://www.log24.com/log/pix11/110420-DarkAndStormy-Logicomix.jpg"
          },
          {
            "url": "https://scontent-atl3-1.cdninstagram.com/vp/f2a627dec636fcfaff011a30221ff6cc/5BB206C0/t51.2885-15/e35/26866491_214283009139346_7442216113854742528_n.jpg"
          },
          {
            "url": "https://s1-ssl.dmcdn.net/iOsV4/x240-M1-.jpg"
          },
          {
            "url": "https://i.pinimg.com/236x/ac/64/6b/ac646be2e343ff5107032e8eaf833d93--sterling-publishing-letterhead.jpg"
          },
          {
            "url": "http://3.bp.blogspot.com/-I8v5r2JASKw/UFirSS82n9I/AAAAAAAAAsg/vadLO5LLmrI/s1600/IMG_0292.PNG"
          }
        ],
        "bestGuessLabels": [
          {
            "label": "logicomix [book]",
            "languageCode": "en"
          }
        ]
      }
    }
  ]
}`
