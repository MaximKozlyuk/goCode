package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

var (
	wb bank
)

//  wordBank.xlsx

func init() {

	fmt.Println(os.Getwd())
	scanExelFile()
}

func main() {

	startQuiz(10)
	fmt.Println("Done!")
}

// рандомим 10 слов для квиза, варианты рандомим из общей массы
// todo в конце добавить выбор, повторить ли слова с неверными ответами
func startQuiz(wordsAmount int) error {
	if (wordsAmount < 0) || (wordsAmount > len(wb.words)) {
		return errors.New("amount of words for quiz invalid")
	}
	if &wb == nil {
		return errors.New("word bank is empty, try ro load it")
	}
	var randId, answer, wordCount int
	quizWords := make(map[int]wordsPair)
	for len(quizWords) < wordsAmount {
		rand.Seed(time.Now().UnixNano())
		randId = rand.Intn(len(wb.words))
		quizWords[randId] = wb.words[randId]
	}

	for i := range quizWords {
		fmt.Println(quizWords[i].eng, ":")
		randId = rand.Intn(5)
		for j := 0; j < 5; j++ {
			if j == randId {
				fmt.Println(j+1, ")", quizWords[i].ru)
				wordCount++
			} else {
				fmt.Println(j+1, ")", wb.words[getUniqWordId(&quizWords)].ru)
			}
		}
		go sayWord(quizWords[i].eng)
		fmt.Scanf("%d\n", &answer)
		if answer-1 == randId {
			fmt.Println("True!")
		} else {
			fmt.Println("False!")
			fmt.Println(quizWords[i].eng, " : ", quizWords[i].ru)
		}
	}
	return nil
}

func getUniqWordId(quizWords *map[int]wordsPair) int {
	var id int
	for {
		id = rand.Intn(len(wb.words))
		if contains(quizWords, id) {
			continue
		} else {
			return id
		}
	}
}

func contains(quizWords *map[int]wordsPair, id int) bool {
	for i := range *quizWords {
		if i == id {
			return true
		}
	}
	return false
}

// todo add scan is exists bank and add new wb
func scanExelFile() {
	var (
		fileName string
		err      error
	)
	fmt.Println("enter the exel file name with google translate phrasebook:")
	for {
		_, err = fmt.Scanf("%s", &fileName)
		fmt.Println()
		if err == nil {
			break
		} else {
			fmt.Println("Some error occured: ", err)
		}
	}
	xlsFile, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println("Some error occured when tryed to read exel file: ", err)
	}
	var rows = xlsFile.GetRows(xlsFile.GetSheetName(1))
	wb.words = make([]wordsPair, 0)

	for i := 0; i < len(rows); i++ {
		if rows[i][0] == "английский" {
			wb.addNewWord(&rows[i][2], &rows[i][3])
		} else {
			wb.addNewWord(&rows[i][3], &rows[i][2])
		}
	}
}

type bank struct {
	words    []wordsPair
	lastSave int64 // date not used now
}

type wordsPair struct {
	eng      string
	ru       string
	attempts int
	rightAns int
}

func (b bank) printWords() {
	fmt.Println("PhraseBook contains ", len(b.words), " words:")
	for i := 0; i < len(b.words); i++ {
		fmt.Println(b.words[i])
	}
}

func saveWordBank() {
	//os.Mkdir("./bank/", os.FileMode(0777))
	//os.Chdir("./bank/")
	//writeGob("./amount.gob", wb.amount)
	//writeGob("./eng.gob", &wb.eng)
	//writeGob("./ru.gob", &wb.ru)
	//writeGob("./attempts.gob", &wb.attempts)
	//writeGob("./trueAnsCount.gob", &wb.trueAnsCount)
}

func loadWordBank() {
	//var newBank bank
	//os.Chdir("./bank/")
	//readGob("./amount.gob", &newBank.amount)
	//readGob("./eng.gob", &newBank.eng)
	//readGob("./ru.gob", &newBank.ru)
	//readGob("./attempts.gob", &newBank.attempts)
	//readGob("./trueAnsCount.gob", &newBank.trueAnsCount)
	//wb = newBank
}

func writeGob(filePath string, o interface{}) error {
	file, err := os.Create(filePath)
	defer file.Close()
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(o)
	}
	return err
}

func readGob(filePath string, o interface{}) error {
	file, err := os.Open(filePath)
	defer file.Close()
	if err == nil {
		encoder := gob.NewDecoder(file)
		encoder.Decode(o)
	}
	return err
}

func (b *bank) addNewWord(engWord, ruWord *string) {
	wp := wordsPair{
		*engWord,
		*ruWord,
		0,
		0,
	}
	b.words = append(b.words, wp)
}

func sayWord(w string) {

	// todo перенаправить вывод куда-нибудь в жопу

	var cmd = "say 123"
	cmd += w
	fmt.Println(cmd)
	_, err := exec.Command("bash", "-c", w).Output()
	if err != nil {
		fmt.Println("some shit occured when tryed to talk", err)
	}
}
