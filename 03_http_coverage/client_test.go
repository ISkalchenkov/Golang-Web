package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

const dataset = "dataset.xml"

type TestCase struct {
	Request     SearchRequest
	Response    *SearchResponse
	AccessToken string
	IsError     bool
	ErrorMsg    string
}

// Тест проверяет функциональность клиента.
// дает гарантию, что параметры запроса формируют правильный ответ - слайс пользователей
func TestFunctionality(t *testing.T) {
	cases := []TestCase{
		{
			// Кейс с упорядочиванием "как есть", поиск по строке "Dillard", всего 2 записи, лимит = 10, смещение = 0
			// Проверка поиска по полю Name
			Request: SearchRequest{
				Limit:      10,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     3,
						Name:   "Everett Dillard",
						Age:    27,
						About:  "Sint eu id sint irure officia amet cillum. Amet consectetur enim mollit culpa laborum ipsum adipisicing est laboris. Adipisicing fugiat esse dolore aliquip quis laborum aliquip dolore. Pariatur do elit eu nostrud occaecat.\n",
						Gender: "male",
					},
					{
						ID:     17,
						Name:   "Dillard Mccoy",
						Age:    36,
						About:  "Laborum voluptate sit ipsum tempor dolore. Adipisicing reprehenderit minim aliqua est. Consectetur enim deserunt incididunt elit non consectetur nisi esse ut dolore officia do ipsum.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			AccessToken: serverAccessToken,
			IsError:     false,
		},
		{
			// Кейс с упорядочиванием "как есть", поиск по строке "laborum quis eu consequat", всего 1 записи, лимит = 10, смещение = 0
			// Проверка поиска по полю About
			Request: SearchRequest{
				Limit:      10,
				Offset:     0,
				Query:      "laborum quis eu consequat",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     31,
						Name:   "Palmer Scott",
						Age:    37,
						About:  "Elit fugiat commodo laborum quis eu consequat. In velit magna sit fugiat non proident ipsum tempor eu. Consectetur exercitation labore eiusmod occaecat adipisicing irure consequat fugiat ullamco aliquip nostrud anim irure enim. Duis do amet cillum eiusmod eu sunt. Minim minim sunt sit sit enim velit sint tempor enim sint aliquip voluptate reprehenderit officia. Voluptate magna sit consequat adipisicing ut eu qui.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			AccessToken: serverAccessToken,
			IsError:     false,
		},
		{
			// Кейс с упорядочиванием по убыванию id, query нет, всего 35 записей, лимит = 2, смещение = 2
			Request: SearchRequest{
				Limit:      2,
				Offset:     2,
				Query:      "",
				OrderField: "id",
				OrderBy:    OrderByDesc,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     32,
						Name:   "Christy Knapp",
						Age:    40,
						About:  "Incididunt culpa dolore laborum cupidatat consequat. Aliquip cupidatat pariatur sit consectetur laboris labore anim labore. Est sint ut ipsum dolor ipsum nisi tempor in tempor aliqua. Aliquip labore cillum est consequat anim officia non reprehenderit ex duis elit. Amet aliqua eu ad velit incididunt ad ut magna. Culpa dolore qui anim consequat commodo aute.\n",
						Gender: "female",
					},
					{
						ID:     31,
						Name:   "Palmer Scott",
						Age:    37,
						About:  "Elit fugiat commodo laborum quis eu consequat. In velit magna sit fugiat non proident ipsum tempor eu. Consectetur exercitation labore eiusmod occaecat adipisicing irure consequat fugiat ullamco aliquip nostrud anim irure enim. Duis do amet cillum eiusmod eu sunt. Minim minim sunt sit sit enim velit sint tempor enim sint aliquip voluptate reprehenderit officia. Voluptate magna sit consequat adipisicing ut eu qui.\n",
						Gender: "male",
					},
				},
				NextPage: true,
			},
			AccessToken: serverAccessToken,
			IsError:     false,
		},
		{
			// Кейс с упорядочиванием по возрастанию Name, query нет, всего 35 записей, лимит = 1, смещение = 0
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "",
				OrderField: "name",
				OrderBy:    OrderByAsc,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     15,
						Name:   "Allison Valdez",
						Age:    21,
						About:  "Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n",
						Gender: "male",
					},
				},
				NextPage: true,
			},
			AccessToken: serverAccessToken,
			IsError:     false,
		},
		{
			// Кейс с упорядочиванием по возрастанию Name, query нет, всего 35 записей, лимит = 1, смещение = 0
			// Кейс аналогичный тому, что расположен выше. Удостоверяемся, что OrderField нечувствителен к регистру
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "",
				OrderField: "NaMe",
				OrderBy:    OrderByAsc,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     15,
						Name:   "Allison Valdez",
						Age:    21,
						About:  "Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n",
						Gender: "male",
					},
				},
				NextPage: true,
			},
			AccessToken: serverAccessToken,
			IsError:     false,
		},
		{
			// Кейс с упорядочиванием по возрастанию Name, query нет, всего 35 записей, лимит = 1, смещение = 0
			// Кейс аналогичный тому, что расположен выше. Проверяем, что по умолчанию упорядочивает по Name
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsc,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     15,
						Name:   "Allison Valdez",
						Age:    21,
						About:  "Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n",
						Gender: "male",
					},
				},
				NextPage: true,
			},
			AccessToken: serverAccessToken,
			IsError:     false,
		},
		{
			// Кейс с упорядочиванием по убыванию Age, query нет, всего 35 записей, лимит = 1, смещение = 0
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "",
				OrderField: "age",
				OrderBy:    OrderByDesc,
			},
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     13,
						Name:   "Whitley Davidson",
						Age:    40,
						About:  "Consectetur dolore anim veniam aliqua deserunt officia eu. Et ullamco commodo ad officia duis ex incididunt proident consequat nostrud proident quis tempor. Sunt magna ad excepteur eu sint aliqua eiusmod deserunt proident. Do labore est dolore voluptate ullamco est dolore excepteur magna duis quis. Quis laborum deserunt ipsum velit occaecat est laborum enim aute. Officia dolore sit voluptate quis mollit veniam. Laborum nisi ullamco nisi sit nulla cillum et id nisi.\n",
						Gender: "male",
					},
				},
				NextPage: true,
			},
			AccessToken: serverAccessToken,
			IsError:     false,
		},
		{
			// Кейс со смещением превышающим количество выбираемых записей
			Request: SearchRequest{
				Limit:      10,
				Offset:     1000,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			Response: &SearchResponse{
				Users:    []User{},
				NextPage: false,
			},
			AccessToken: serverAccessToken,
			IsError:     false,
		},
	}

	filename = dataset
	PrepareDataAndCheckError(t)
	ts := httptest.NewServer(http.HandlerFunc(ServerSearch))
	for caseNum, item := range cases {
		sc := &SearchClient{
			URL:         ts.URL,
			AccessToken: item.AccessToken,
		}
		response, err := sc.FindUsers(item.Request)

		if err != nil && !item.IsError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if !reflect.DeepEqual(response, item.Response) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, item.Response, response)
		}
	}
	ts.Close()
}

// Тест проверяет, что запрос с превышающим ограничение лимитом, вернет установленное количество записей
func TestLimit(t *testing.T) {
	cases := []TestCase{
		{
			// Кейс со лимитом превышающим ограничение 25
			Request: SearchRequest{
				Limit: 30,
			},
			Response: &SearchResponse{
				Users:    make([]User, 25), // в данном тесте важна только длина полученного слайса
				NextPage: true,
			},
			AccessToken: serverAccessToken,
			IsError:     false,
		},
	}

	filename = dataset
	PrepareDataAndCheckError(t)
	ts := httptest.NewServer(http.HandlerFunc(ServerSearch))
	for caseNum, item := range cases {
		sc := &SearchClient{
			URL:         ts.URL,
			AccessToken: item.AccessToken,
		}
		response, err := sc.FindUsers(item.Request)

		if err != nil && !item.IsError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if len(item.Response.Users) != len(response.Users) {
			t.Errorf("[%d] wrong number of selected users, expected %d, got %d", caseNum, len(item.Response.Users), len(response.Users))
		}
		if item.Response.NextPage != response.NextPage {
			t.Errorf("[%d] wrong flag NextPage, expected %t, got %t", caseNum, item.Response.NextPage, response.NextPage)
		}
	}
	ts.Close()
}

// Тест проверяет возникновение ошибок, связанных с некорректным составлением запроса:
// Limit, Offset, OrderField, OrderBy, AccessToken
func TestRequestErrors(t *testing.T) {
	cases := []TestCase{
		{
			// Кейс с OrderField = Gender. Невалидный OrderField
			Request: SearchRequest{
				Limit:      10,
				Offset:     0,
				Query:      "",
				OrderField: "Gender",
				OrderBy:    OrderByAsIs,
			},
			Response:    nil,
			AccessToken: serverAccessToken,
			IsError:     true,
			ErrorMsg:    "OrderFeld Gender invalid",
		},
		{
			// Кейс с OrderBy = 999. Невалидный OrderBy
			Request: SearchRequest{
				Limit:      10,
				Offset:     0,
				Query:      "",
				OrderField: "",
				OrderBy:    999,
			},
			Response:    nil,
			AccessToken: serverAccessToken,
			IsError:     true,
			ErrorMsg:    "unknown bad request error: OrderBy invalid",
		},
		{
			// Кейс с Limit = -1. Невалидный Limit
			Request: SearchRequest{
				Limit:      -1,
				Offset:     0,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			Response:    nil,
			AccessToken: serverAccessToken,
			IsError:     true,
			ErrorMsg:    "limit must be > 0",
		},
		{
			// Кейс с Offset = -1. Невалидный Offset
			Request: SearchRequest{
				Limit:      10,
				Offset:     -1,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			Response:    nil,
			AccessToken: serverAccessToken,
			IsError:     true,
			ErrorMsg:    "offset must be > 0",
		},
		{
			// Кейс с AccessToken = "WrongAccessToken". Неверный AccessToken
			Request: SearchRequest{
				Limit:      10,
				Offset:     0,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsIs,
			},
			Response:    nil,
			AccessToken: "WrongAccessToken",
			IsError:     true,
			ErrorMsg:    "bad AccessToken",
		},
	}

	filename = dataset
	PrepareDataAndCheckError(t)
	ts := httptest.NewServer(http.HandlerFunc(ServerSearch))
	CheckErrorCases(t, ts.URL, cases)
	ts.Close()
}

// Тест симулирует ошибку чтения(открытия файла) данных о пользователях на сервере
// Возвращает StatusInternalServerError клиенту
func TestServerReadFileError(t *testing.T) {
	cases := []TestCase{
		{
			Request:     SearchRequest{},
			Response:    nil,
			AccessToken: serverAccessToken,
			IsError:     true,
			ErrorMsg:    "SearchServer fatal error",
		},
	}

	filename = "nonexistent-file.xml"
	usersData = Users{} // Зануляем данные, если были считаны ранее
	PrepareUsersData()  //nolint:errcheck
	ts := httptest.NewServer(http.HandlerFunc(ServerSearch))
	CheckErrorCases(t, ts.URL, cases)
	ts.Close()
}

// Тест симулирует ошибку десериализации данных из файла с информацией о пользователях
// Возвращает StatusInternalServerError клиенту
func TestServerUnmarshalFileError(t *testing.T) {
	cases := []TestCase{
		{
			Request:     SearchRequest{},
			Response:    nil,
			AccessToken: serverAccessToken,
			IsError:     true,
			ErrorMsg:    "SearchServer fatal error",
		},
	}

	filename = "corrupted-file.xml"
	usersData = Users{}
	corruptedUsersData := []byte("[]corrupted{};;; ]users data[]")
	err := os.WriteFile(filename, corruptedUsersData, 0666)
	if err != nil {
		t.Errorf("file write error: %v", err)
	}

	PrepareUsersData() //nolint:errcheck
	ts := httptest.NewServer(http.HandlerFunc(ServerSearch))
	CheckErrorCases(t, ts.URL, cases)
	os.Remove(filename)
	ts.Close()
}

// Тест проверяет обработку клиентом ошибки, возникающей в результате выполнения Unmarshal сообщения при ответе со статусом BadRequest
// Тело ответа невалидно
func TestBadRequestUnmarshalError(t *testing.T) {
	cases := []TestCase{
		{
			Request:     SearchRequest{},
			Response:    nil,
			AccessToken: serverAccessToken,
			IsError:     true,
			ErrorMsg:    "cant unpack error json: invalid character '[' after top-level value",
		},
	}

	badJSONHandler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error": "bad json response"}[]`, http.StatusBadRequest)
	}

	filename = dataset
	PrepareDataAndCheckError(t)
	ts := httptest.NewServer(http.HandlerFunc(badJSONHandler))
	CheckErrorCases(t, ts.URL, cases)
	ts.Close()
}

// Тест проверяет обработку клиентом ошибки, возникающей в результате выполнения Unmarshal сообщения при успешном ответе
// Тело ответа невалидно
func TestSuccessRequestUnmarshalError(t *testing.T) {
	cases := []TestCase{
		{
			Request:     SearchRequest{},
			Response:    nil,
			AccessToken: serverAccessToken,
			IsError:     true,
			ErrorMsg:    "cant unpack result json: invalid character '[' after object key:value pair",
		},
	}

	badJSONHandler := func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`{"error": "bad json response"[]`))
		if err != nil {
			log.Printf("error writing data to the connection: %v", err)
		}
	}

	filename = dataset
	PrepareDataAndCheckError(t)
	ts := httptest.NewServer(http.HandlerFunc(badJSONHandler))
	CheckErrorCases(t, ts.URL, cases)
	ts.Close()
}

// Тест проверяет обработку клиентом ошибки при таймауте выполнения запроса
func TestTimeoutError(t *testing.T) {
	cases := []TestCase{
		{
			Request:     SearchRequest{},
			Response:    nil,
			AccessToken: serverAccessToken,
			IsError:     true,
			ErrorMsg:    "timeout for limit=1&offset=0&order_by=0&order_field=&query=",
		},
	}
	timeoutHandler := func(w http.ResponseWriter, r *http.Request) {
		sleepTime := client.Timeout + time.Millisecond*100
		time.Sleep(sleepTime)
	}

	filename = dataset
	PrepareDataAndCheckError(t)
	ts := httptest.NewServer(http.HandlerFunc(timeoutHandler))
	CheckErrorCases(t, ts.URL, cases)
	ts.Close()
}

// Тест проверяет обработку клиентом прочих ошибок при выполнении запроса
func TestClientDoErrors(t *testing.T) {
	cases := []TestCase{
		{
			Request:     SearchRequest{},
			Response:    nil,
			AccessToken: serverAccessToken,
			IsError:     true,
			ErrorMsg:    `unknown error Get "?limit=1&offset=0&order_by=0&order_field=&query=": unsupported protocol scheme ""`,
		},
	}

	filename = dataset
	PrepareDataAndCheckError(t)
	ts := httptest.NewServer(http.HandlerFunc(ServerSearch))
	CheckErrorCases(t, "", cases) // URL = ""
	ts.Close()
}

func PrepareDataAndCheckError(t *testing.T) {
	err := PrepareUsersData()
	if err != nil {
		t.Errorf("failed to prepare users data: %v", err)
	}
}

func CheckErrorCases(t *testing.T, url string, cases []TestCase) {
	for caseNum, item := range cases {
		sc := &SearchClient{
			URL:         url,
			AccessToken: item.AccessToken,
		}
		response, err := sc.FindUsers(item.Request)

		if err == nil && item.IsError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}

		if err.Error() != item.ErrorMsg {
			t.Errorf("[%d] wrong error message, expected %s, got %s",
				caseNum, item.ErrorMsg, err.Error())
		}

		if !reflect.DeepEqual(response, item.Response) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v",
				caseNum, item.Response, response)
		}
	}
}

// Ниже расположен тест покрытия кода сервера, который невозможно покрыть путем написания тестов для клиента
// поскольку клиент проверяет лимиты и смещения перед совершением запроса,
// невалидные лимиты и смещения передаться не могут, а проверка на сервере имеется

type ServerTestCase struct {
	AccessToken string
	Request     SearchRequest
	Response    string
	StatusCode  int
}

// Тест с недостающим покрытием сервера
func TestServerSearch(t *testing.T) {
	cases := []ServerTestCase{
		// Кейс с невалидным лимитом
		{
			AccessToken: serverAccessToken,
			Request: SearchRequest{
				Limit: -1,
			},
			Response:   `{"error": "Limit invalid"}` + "\n",
			StatusCode: 400,
		},
		// Кейс с невалидным смещением
		{
			AccessToken: serverAccessToken,
			Request: SearchRequest{
				Limit:  1,
				Offset: -1,
			},
			Response:   `{"error": "Offset invalid"}` + "\n",
			StatusCode: 400,
		},
	}
	filename = dataset
	PrepareDataAndCheckError(t)
	for caseNum, item := range cases {
		searcherParams := url.Values{}
		searcherParams.Add("limit", strconv.Itoa(item.Request.Limit))
		searcherParams.Add("offset", strconv.Itoa(item.Request.Offset))
		searcherParams.Add("query", item.Request.Query)
		searcherParams.Add("order_field", item.Request.OrderField)
		searcherParams.Add("order_by", strconv.Itoa(item.Request.OrderBy))

		url := "http://example.com/api/search?" + searcherParams.Encode()
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Add("AccessToken", item.AccessToken)
		w := httptest.NewRecorder()

		ServerSearch(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("[%d] failed to read body: %v", caseNum, err)
		}

		bodyStr := string(body)
		if bodyStr != item.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				caseNum, bodyStr, item.Response)
		}
	}
}
