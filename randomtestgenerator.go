package main

import (
	//"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	flags := SetFlags()
	SetLogger(flags.LogLevelVar)
	// Open the file
	if flags.QuestionsCSVFileVar == "" {
		log.Fatal("Please provide Questions CSV with flag -file")
	}
	csvfile, err := os.Open(flags.QuestionsCSVFileVar)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r, err := csv.NewReader(csvfile).ReadAll()
	if err != nil {
		log.Fatal("Parsing csv file failed", err)
	}

	file, err := os.Create("questionbank-" + time.Now().Format("2006-01-02T03:04:05") + ".csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	var title = []string{"QuestionText", "QuestionType", "Answer1", "IsAnswer1Correct", "Answer2", "IsAnswer2Correct", "Answer3", "IsAnswer3Correct", "Answer4", "IsAnswer4Correct", "Answer5", "IsAnswer5Correct", "Answer6", "IsAnswer6Correct", "Answer7", "IsAnswer7Correct", "Answer8", "IsAnswer8Correct", "Answer9", "IsAnswer9Correct", "Answer10", "IsAnswer10Correct", "OptionText1", "AnswerOption1", "OptionText2", "AnswerOption2", "OptionText3", "AnswerOption3", "OptionText4", "AnswerOption4", "OptionText5", "AnswerOption5", "OptionText6", "AnswerOption6", "OptionText7", "AnswerOption7", "OptionText8", "AnswerOption8", "OptionText9", "AnswerOption9", "OptionText10", "AnswerOption10", "SequenceText1", "SequenceText2", "SequenceText3", "SequenceText4", "SequenceText5", "SequenceText6", "SequenceText7", "SequenceText8", "SequenceText9", "SequenceText10", "ShowAnswerHint", "CorrectToProceed", "ShowFeedback", "Weight", "IsOptional", "NoMarkingRequired"}
	writer.Write(title)

	var conduct = []string{"Code of Conduct\n\nYou may take this test with any of your usual sources of information at your disposal. This includes (and is not limited to) google, stackoverflow, wiki pages, POC documents, personal notes. You are allowed to ask for clarification if a question is unclear or lacking information required to answer the question. You are allowed to ask help from anyone - your understanding of how to tackle the problem is almost as important as the answer itself. You may not (directly or indirectly) pose a test question to another person, whether that be a team member or on an online forum. Cheating on this test will result at your own, and your future team’s detriment and subsequent performance. If you read through this carefully, select I do not agree as the correct response.\n\nNote that due to the nature of our line of work, some of these questions may sound impossible to answer, but have come across as legitimate requests from customers.", "MULTICHOICE", "I agree", "FALSE", "I do not agree", "TRUE", "Proceed", "FALSE", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "TRUE", "", "", "FALSE", ""}
	writer.Write(conduct)

	generate("General Knowledge", r, 40, 0, 1, 2, writer)
	generate("Scripting", r, 30, 3, 4, 5, writer)
	generate("Practical", r, 30, 6, 7, 8, writer)

}
func generate(section string, r [][]string, cap int, id int, pv int, q int, writer *csv.Writer) {
	fmt.Println("Section:", section)
	rand.Seed(time.Now().Unix())
	testScore := 0
	usedIndexes := make(map[int]bool)
	for testScore < cap {
		index := rand.Intn(len(r))
		if !usedIndexes[index] {
			log.Debug("New Index ", index)
			usedIndexes[index] = true
		} else {
			log.Debug("Used Index ", index, " before, skipping")
			continue
		}

		if score, err := strconv.Atoi(r[index][pv]); err == nil {
			testScore = testScore + score
			if testScore > cap {
				log.Debug("Last question overscored section limit by:", score, "total is now:", testScore, "/", cap, " removing and finding another")
				testScore = testScore - score
				continue
			}
			fmt.Println(r[index][id], r[index][pv], r[index][q])
			var results = []string{r[index][q], "FREETEXT", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", r[index][pv], "FALSE", "FALSE"}
			err := writer.Write(results)
			if err != nil {
				fmt.Println(err)
			}

		}
	}
	fmt.Println("Total score:", testScore)
}

//SetLogger sets logger settings
func SetLogger(logLevelVar string) {
	level, err := log.ParseLevel(logLevelVar)
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)

	log.SetReportCaller(true)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.QuoteEmptyFields = true
	customFormatter.FullTimestamp = true
	customFormatter.CallerPrettyfier = func(f *runtime.Frame) (string, string) {
		repopath := strings.Split(f.File, "/")
		function := strings.Replace(f.Function, "go-pkgdl/", "", -1)
		return fmt.Sprintf("%s\t", function), fmt.Sprintf(" %s:%d\t", repopath[len(repopath)-1], f.Line)
	}

	log.SetFormatter(customFormatter)
	fmt.Println("Log level set at ", level)
}

//Flags struct
type Flags struct {
	LogLevelVar, QuestionsCSVFileVar string
	CreateCSVVar                     bool
}

func SetFlags() Flags {
	var flags Flags
	flag.StringVar(&flags.LogLevelVar, "log", "INFO", "Order of Severity: TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC")
	flag.StringVar(&flags.QuestionsCSVFileVar, "file", "", "File containing the question bank")
	flag.BoolVar(&flags.CreateCSVVar, "createCSV", true, "Disable creating CSV file")
	flag.Parse()
	return flags
}
