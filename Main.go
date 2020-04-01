package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/go-gota/gota/dataframe"

	"github.com/adshao/go-binance"
)

var (
	apiKey    = "entrer votre API"
	secretKey = "entrer votre API"
)

/// Pour retourner la data de binance ic
func get_data(b string) []*binance.Kline { // on récupere ce type là qui est un ptr slice
	client := binance.NewClient(apiKey, secretKey)
	data, _ := client.NewKlinesService().Symbol("BTCUSDT").
		Interval(b).Limit(1000).Do(context.Background())
	Start := data[0].OpenTime - 6 /// changer le 6 selon le tmps 6=60s 360=1h etc...
	lendata := len(data)
	for i := 0; i < 70; i++ { ///// Changer le 70 pour la taille de la data ici c'est 71000 éléments de 1min
		datab, _ := client.NewKlinesService().Symbol("BTCUSDT").Interval(b).Limit(1000).EndTime(Start).Do(context.Background())
		for i := 0; i < lendata; i++ {
			data = append(data, datab[i])
		}
		Start = datab[0].OpenTime - 6
	}
	sort.SliceStable(data, func(i, j int) bool { //Fonction de tri
		return data[i].OpenTime < data[j].OpenTime
	})
	fmt.Println(len(data))
	return data
}

//////////////////Passer en slice struc string pour ensuite passer en dataframe pour étudier.
func converstruct(klines []*binance.Kline) interface{} {
	type dat struct {
		OpenTime  string `json:"openTime"`
		Open      string `json:"open"`
		High      string `json:"high"`
		Low       string `json:"low"`
		Close     string `json:"close"`
		Volume    string `json:"volume"`
		CloseTime string `json:"closeTime"`
	}

	data := make([]dat, len(klines))
	test := 600

	for k := 0; k < test; k++ {
		s := strconv.FormatInt(klines[k].OpenTime, 10)
		d := strconv.FormatInt(klines[k].CloseTime, 10)
		data[k] = dat{
			OpenTime:  s,
			Open:      klines[k].Open,
			High:      klines[k].High,
			Low:       klines[k].Low,
			Close:     klines[k].Close,
			Volume:    klines[k].Volume,
			CloseTime: d,
		}
	}
	df := dataframe.LoadStructs(data)
	f, _ := os.Create("Test.CSV")
	df.WriteCSV(f)

	return data
}

//////Fonction d'affichage d'un Tableau 
func printtab(tab []int64) {
	for _, k := range tab {
		fmt.Println(k)
	}
	fmt.Println("Taille est", len(tab))
}

/////Fonction backtest pour le type []*binance.kline
//// tab doit être un tableau avec les times d'achats
//// tab2 doit être un tableau avec les times de ventes
func backtest1(tab []int64, tab2 []int64, klines []*binance.Kline) {
	lentab := len(tab)
	var tabretour []int64
	a := len(klines)
	wa := 10000.000
	for i := 0; i < lentab; i++ {
		k2 := tab[i]
		v := tab2[i]
		for i := 0; i < a; i++ {
			k1 := klines[i].CloseTime
			if k1 == k2 {
				for g := 0; g < a; g++ {
					v2 := klines[g].CloseTime
					if v == v2 {
						bec := klines[g].Close
						becf, _ := strconv.ParseFloat(bec, 64)
						be := klines[i].Close
						bef, _ := strconv.ParseFloat(be, 64)
						pnlf := (becf - bef) / bef
						wa = wa + pnlf*wa
						if becf > bef { //On regarde combien de fois on a gagné
							tabretour = append(tabretour, 1)
							break
						} else {
							tabretour = append(tabretour, 0)
							break
						}
					}
				}
				break
			}
		}
	}
	a1 := 0
	b1 := len(tabretour)
	fmt.Println("Nombre d'élèments", b1)
	for _, k := range tabretour {
		//fmt.Println(k)
		if k == 1 {
			a1 = a1 + 1
			continue
		}
	}
	a2 := float32(a1)
	b2 := float32(b1)
	res := (a2 / b2) * 100
	fmt.Println("le pourcentage de réussite ", res)
	fmt.Println("Notre portefeuille test", wa)
	g := klines[a-1].Close
	gf, _ := strconv.ParseFloat(g, 32)
	n := klines[0].Close
	gi, _ := strconv.ParseFloat(n, 32)
	wa2 := 10000.000
	rt := (gf - gi) / gi
	wa2 = wa2 + rt*wa2
	fmt.Println("Le portfeuille sans bot", wa2) 
	//Cela nous montre si on avait passer un ordre d'achat au dévut de la data et de vente à la fin
}

//// Fonction pour convertir la Data en 0 et 1, notamment pour du machine learning
func Dataconvert(b string) []int64 {
	klines := get_data(b)
	a := len(klines)
	var tabtran []int64
	for i := 0; i < a-1; i++ {
		d := klines[i+1].Close
		d1, _ := strconv.ParseFloat(d, 64)
		c := klines[i].Close
		c1, _ := strconv.ParseFloat(c, 64)
		if d1 > c1 {
			tabtran = append(tabtran, 1)
			continue
		} else {
			tabtran = append(tabtran, 0)
			continue
		}
	}
	return tabtran
}
///Algo test pour tester nos fonctions, le but est de chercher deux klines de suite égaux sur certains points
func Lesklinesegaux(b string, klines []*binance.Kline) { //Petit algo pour des klines égaux
	var tabachat []int64
	var tabvente []int64
	////Fonction pour des bougies égaux
	a := len(klines)
	for i := 0; i < a-1; i++ {
		d := klines[i].Low
		c := klines[i].Close
		c1 := klines[i+1].Open
		if d == c && c == c1 {
			g := klines[i].CloseTime
			d := klines[i+1].CloseTime
			tabachat = append(tabachat, g)
			tabvente = append(tabvente, d)
			continue
		} else {
			continue
		}

	}
///On test sur notre backtest
	backtest1(tabachat, tabvente, klines)

}

///// Fonction main////////////
func main() {
	
	klines := get_data("1m") // On recupère la data en 1m
	data := converstruct(klines) // on convertie la data 
	df := dataframe.LoadStructs(data) // on convertie en dataframe
	fmt.Println(df) // on affiche
	Lesklinesegaux("1m", klines) // on applique notre algo
}
