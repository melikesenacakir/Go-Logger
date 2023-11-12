package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type User struct {
	username string
	password string
}

func savelog(dbname string, status bool) {
	logDosyam, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Dosya açılamadı")
		os.Exit(0)
	} else {
		log.SetFlags(0)
		defer logDosyam.Close()
		log.SetOutput(logDosyam)
		log.Printf("Kullanıcı adı: " + dbname)
		log.Printf("Giriş Tarihi ve saati: %s", time.Now().Format("2006-01-02 15:04:05"))

		if status == true {
			log.Println("Giriş Durumu: " + "\033[32m" + "Başarılı" + "\033[0m")
		} else {
			log.Println("Giriş Durumu: " + "\033[31m" + "Başarısız" + "\033[0m")
		}
		log.Println("***************************")
	}
}

func controlLogin(dbname string, dbpass string, username string, pass string) bool {
	if dbname == username {
		if dbpass == pass {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func main() {
	logintype := 0
	var username string
	var pass string
	var exittype int

	admin := User{
		username: "admin",
		password: "123",
	}
	student := User{
		username: "s1",
		password: "s12",
	}

	burstyLimiter := make(chan time.Time, 3)
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}
	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	burstyRequests := make(chan int, 5)
	for i := 4; i >= 0; i-- {
		burstyRequests <- i
	}
	close(burstyRequests)

	fmt.Println("Lütfen giriş türünüzü seçin:\n0 - Admin Girişi\n1 - Öğrenci Girişi")
	fmt.Scan(&logintype)

	if logintype == 0 {
		res := false
		for req := range burstyRequests {
			<-burstyLimiter
			fmt.Print("\nAdmin kullanıcı adı: ")
			fmt.Scan(&username)
			fmt.Print("Admin şifre: ")
			fmt.Scan(&pass)
			res = controlLogin(admin.username, admin.password, username, pass)
			savelog(username, res)
			if res != true && req != 0 {
				fmt.Println("\033[33m" + "Geçersiz kullanıcı adı veya şifre! kalan deneme hakkınız " + fmt.Sprintf("%d", req) + "\033[0m")
			} else {
				break
			}
		}
		if res == true {
			exittype = 0
			fmt.Println("\033[32m" + "Başarılı Giriş\n*********************" + "\033[0m")
			for exittype != 1 {
				fmt.Println("Lütfen bir işlem seçin:\n0 - Logları Görüntüle\n1 - Çıkış Yap")
				fmt.Scan(&exittype)
				if exittype == 0 {
					logDosyam, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE, 0644)
					if err != nil {
						fmt.Println("Dosya açılamadı")
						os.Exit(0)
					} else {
						defer logDosyam.Close()
						scanner := bufio.NewScanner(logDosyam)
						for scanner.Scan() {
							line := scanner.Text()
							fmt.Println(line)
						}
					}
				}
			}
			fmt.Println("\033[31m" + "Çıkış Yapılmıştır\n*********************" + "\033[0m")
			os.Exit(0)
		}

	} else {
		res := false
		for req := range burstyRequests {
			<-burstyLimiter
			fmt.Print("\nÖğrenci kullanıcı adı: ")
			fmt.Scan(&username)
			fmt.Print("Öğrenci şifre: ")
			fmt.Scan(&pass)
			res = controlLogin(student.username, student.password, username, pass)
			savelog(username, res)
			if res != true && req != 0 {
				fmt.Println("\033[33m" + "Geçersiz kullanıcı adı veya şifre! kalan deneme hakkınız " + fmt.Sprintf("%d", req) + "\033[0m")
			} else {
				break
			}
		}
		if res == true {
			fmt.Println("\033[32m" + "Başarılı Giriş\n*********************" + "\033[0m")
			fmt.Println("Lütfen bir işlem seçin:\n1 - Çıkış Yap")
			fmt.Scan(&exittype)
			if exittype == 1 {
				fmt.Println("\033[31m" + "Çıkış Yapılmıştır\n*********************" + "\033[0m")
				os.Exit(0)
			}

		}
	}
}
