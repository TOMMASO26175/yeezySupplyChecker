package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Products []struct {
	Price          int       `json:"price"`
	ProductID      string    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	ProductModelID string    `json:"product_model_id"`
	PreviewTo      time.Time `json:"previewTo"`
	Image          struct {
		Alt   string `json:"alt"`
		Link  string `json:"link"`
		Title string `json:"title"`
	} `json:"image"`
	GroupSortID     int      `json:"groupSortId"`
	GroupItemSortID int      `json:"groupItemSortId"`
	CalloutMessages []string `json:"calloutMessages"`
	Color           string   `json:"color"`
	PreOrderable    bool     `json:"preOrderable"`
	IsWaitingRoom   bool     `json:"isWaitingRoom"`
}

type Availability struct {
	ID                 string `json:"id"`
	AvailabilityStatus string `json:"availability_status"`
	VariationList      []struct {
		Sku                string `json:"sku"`
		Size               string `json:"size"`
		Availability       int    `json:"availability"`
		AvailabilityStatus string `json:"availability_status"`
	} `json:"variation_list"`
}

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ModelNumber string `json:"model_number"`
	ProductType string `json:"product_type"`
	MetaData    struct {
		PageTitle   string `json:"page_title"`
		SiteName    string `json:"site_name"`
		Description string `json:"description"`
		Keywords    string `json:"keywords"`
		Canonical   string `json:"canonical"`
	} `json:"meta_data"`
	YeezyPDPCallout []string `json:"yeezyPDPCallout"`
	ViewList        []struct {
		Type     string `json:"type"`
		ImageURL string `json:"image_url"`
		Source   string `json:"source"`
	} `json:"view_list"`
	PricingInformation struct {
		StandardPrice      int `json:"standard_price"`
		StandardPriceNoVat int `json:"standard_price_no_vat"`
		CurrentPrice       int `json:"currentPrice"`
	} `json:"pricing_information"`
	AttributeList struct {
		IsWaitingRoomProduct     bool     `json:"isWaitingRoomProduct"`
		BadgeText                string   `json:"badge_text"`
		BadgeStyle               string   `json:"badge_style"`
		Brand                    string   `json:"brand"`
		Collection               []string `json:"collection"`
		Category                 string   `json:"category"`
		Color                    string   `json:"color"`
		ReturnType               string   `json:"return_type"`
		MtbrFlag                 bool     `json:"mtbr_flag"`
		Gender                   string   `json:"gender"`
		Personalizable           bool     `json:"personalizable"`
		MandatoryPersonalization bool     `json:"mandatory_personalization"`
		Customizable             bool     `json:"customizable"`
		Pricebook                string   `json:"pricebook"`
		Sale                     bool     `json:"sale"`
		Outlet                   bool     `json:"outlet"`
		IsCnCRestricted          bool     `json:"isCnCRestricted"`
		SizeChartLink            string   `json:"size_chart_link"`
		Sport                    []string `json:"sport"`
		SizeFitBar               struct {
			Value               string `json:"value"`
			SelectedMarkerIndex int    `json:"selectedMarkerIndex"`
			MarkerCount         int    `json:"markerCount"`
		} `json:"size_fit_bar"`
		PreviewTo         time.Time `json:"preview_to"`
		ComingSoonSignup  bool      `json:"coming_soon_signup"`
		MaxOrderQuantity  int       `json:"max_order_quantity"`
		ProductType       []string  `json:"productType"`
		SearchColor       string    `json:"search_color"`
		SpecialLaunch     bool      `json:"specialLaunch"`
		SpecialLaunchType string    `json:"specialLaunchType"`
		SearchColorRaw    string    `json:"search_color_raw"`
	} `json:"attribute_list"`
	ProductDescription struct {
		Title             string   `json:"title"`
		Usps              []string `json:"usps"`
		DescriptionAssets struct {
		} `json:"description_assets"`
	} `json:"product_description"`
	RecommendationsEnabled bool          `json:"recommendationsEnabled"`
	ProductLinkList        []interface{} `json:"product_link_list"`
}

func main() {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	pid := "FY4567"
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Select what you want to retrieve: ")
	fmt.Println("1. Products ")
	fmt.Println("2. Specific product")
	answ, _ := reader.ReadString('\n')
	answ = strings.TrimRight(answ, "\r\n")
	if answ == "1" {
		req, err := http.NewRequest("GET", "https://www.yeezysupply.com/api/yeezysupply/products/bloom", nil)
		req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36`)
		HandleReq(req, myClient, err, Products{})
	} else if answ == "2" {
		fmt.Println("Select what you want to retrieve:")
		fmt.Println("1. Availability ")
		fmt.Println("2. Product data")
		answ, _ := reader.ReadString('\n')
		answ = strings.TrimRight(answ, "\r\n")
		//Request to get the cookie for accessing to product availability and the product itself
		cookieString := RequestCookie("https://www.yeezysupply.com/product/" + pid)
		name, value := splitCookie(cookieString)
		cookie := &http.Cookie{Name: name, Value: value}
		switch answ {
		case "1":
			req, err := http.NewRequest("GET", "https://www.yeezysupply.com/api/products/"+pid+"/availability", nil)
			req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36`)
			req.AddCookie(cookie)
			HandleReq(req, myClient, err, Availability{})
			break
		case "2":
			req, err := http.NewRequest("GET", "https://www.yeezysupply.com/api/products/"+pid, nil)
			req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36`)
			req.AddCookie(cookie)
			HandleReq(req, myClient, err, Product{})
			break
		}

	} else {
		err := errors.New("wrong imput")
		if err != nil {
			fmt.Println(err)
		}
	}

}

func HandleReq(req *http.Request, client *http.Client, err error, result interface{}) {
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		log.Printf("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("response: %s", body)
	}
	fmt.Println(PrettyPrint(result))
}

func splitCookie(cookie string) (string, string) {
	split := strings.SplitN(cookie, "=", 2)
	return split[0], strings.Split(split[1], ";")[0]
}

func RequestCookie(url string) string {
	var client = &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36`)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("No response from request")
	}
	cookies := resp.Header.Get("Set-Cookie")
	return cookies

}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
