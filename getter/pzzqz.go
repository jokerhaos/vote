package getter

import (
	"strings"

	"github.com/Aiicy/htmlquery"
	"github.com/henson/proxypool/pkg/models"
)

// PZZQZ get ip from http://pzzqz.com/
func PZZQZ() (result []*models.IP) {
	pollURL := "http://pzzqz.com/"
	doc, _ := htmlquery.LoadURL(pollURL)
	trNode, err := htmlquery.Find(doc, "//table[@class='table table-hover']//tbody//tr")
	if err != nil {
		// clog.Warn("pzzqz:", err.Error())
	}
	for i := 0; i < len(trNode); i++ {
		tdNode, _ := htmlquery.Find(trNode[i], "//td")
		ip := htmlquery.InnerText(tdNode[0])
		port := htmlquery.InnerText(tdNode[1])
		Type := htmlquery.InnerText(tdNode[4])

		IP := models.NewIP()
		IP.Data = ip + ":" + port
		IP.Type1 = strings.ToLower(Type)
		IP.Source = "pzzqz"
		result = append(result, IP)
	}

	// clog.Info("[pzzqz] done")
	return
}
