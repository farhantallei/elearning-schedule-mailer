package main

import (
	"elearning-schedule-mailer/services"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"time"
)

func mailer(courseName string, endTime string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	stopAfter := time.After(time.Until(getEndTime(endTime)))
	stopFetch := make(chan bool)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Tugas dijalankan pada: ", time.Now())

			scheduleResponse, err := services.FetchSchedule(courseName)
			if err != nil {
				fmt.Println("Error fetching schedule: ", err)
				continue
			}

			if scheduleResponse.LinkMedia == nil {
				fmt.Println("No media link found.")
				continue
			}

			subject := fmt.Sprintf("ELR LINK KULIAH: %s", courseName)
			body := fmt.Sprintf("Dosen: %s<br>Pertemuan: %s<br>Topik: %s<br>Link kuliah: %s", scheduleResponse.LecturerName, scheduleResponse.CourseTopic, defaultIfNil(scheduleResponse.Noted), defaultIfNil(scheduleResponse.LinkMedia))

			_, err = services.SendMail(subject, body)

			if err != nil {
				fmt.Println("Error sending mail: ", err)
				continue
			}

			fmt.Println("Mail sent successfully.")
			stopFetch <- true
		case <-stopAfter:
			fmt.Println("Interval dihentikan pada: ", time.Now())
			return
		case <-stopFetch:
			fmt.Println("Interval dihentikan pada: ", time.Now())
			return
		}
	}
}

func main() {
	fmt.Println("Starting eLearning schedule mailer...")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Println("Error loading timezone:", err)
		return
	}
	time.Local = loc

	c := cron.New()

	c.AddFunc("45 8 * * 0", func() {
		mailer("Algoritma dan Pemograman 1", "10:30")
	})

	c.AddFunc("50 10 * * 0", func() {
		mailer("Matematika Dasar", "12:35")
	})

	c.AddFunc("0 13 * * 0", func() {
		mailer("Teori Bahasa Otomata", "14:45")
	})

	c.AddFunc("0 15 * * 0", func() {
		mailer("Pemrograman Berorientasi Objek 1", "16:45")
	})

	c.AddFunc("20 20 * * 1", func() {
		mailer("Pemograman Web 1", "22:20")
	})

	c.AddFunc("40 19 * * 3", func() {
		mailer("Struktur Data", "21:40")
	})

	c.AddFunc("0 14 * * 6", func() {
		mailer("Fisika Dasar", "15:45")
	})

	c.AddFunc("15 18 * * 6", func() {
		mailer("Sistem Operasi", "20:00")
	})

	c.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "Hello from Koyeb")
	})

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
