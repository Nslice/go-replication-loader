package argsp

/*
import (
	"fmt"
	"os"
	"encoding/gob"
)
*/
type Student struct {
	Name string
	Age  int32
}

func Ser() {

	/*fmt.Println("Gob Example")
	student := Student{"Ketan Parmar",35}
	err := writeGob("./student.gob",student)
	if err != nil{
		   fmt.Println(err)
	}


	var studentRead = new (Student)
	err = readGob("./student.gob",studentRead)
	if err != nil {
		   fmt.Println(err)
	} else {
		   fmt.Println(studentRead.Name, "\t", studentRead.Age)
	}*/

} /*

func writeGob(filePath string,object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		   encoder := gob.NewEncoder(file)
		   encoder.Encode(object)
	}
	file.Close()
	return err
}

func readGob(filePath string,object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		   decoder := gob.NewDecoder(file)
		   err = decoder.Decode(object)
	}
	file.Close()
	return err
}*/
