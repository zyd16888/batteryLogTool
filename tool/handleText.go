package tool

import (
	"log"
	"os"
	"regexp"
	"strings"
)

func ReText(textline string) (string, bool) {
	re := regexp.MustCompile("(.*)cn.com.imch.mddp.admin.aop.BatteryLogAspect(.*)")
	match := re.FindString(textline)
	if match != "" {
		return match, true
	} else {
		return match, false
	}
}

func WriteFile(textfile string, textline string) error {
	file, err := os.OpenFile(textfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		log.Printf("Cannot open text file: %s, err: [%v]", textfile, err)
		return err
	}
	defer file.Close()

	count, err := file.WriteString(textline + "\n")
	if err != nil {
		log.Printf("Cannot write text line: %s, err: [%v]", textline, err)
		return err
	}
	log.Printf("write count size: %v", count)

	return nil
}

func ReJsonText(text string) string {
	compileRegex := regexp.MustCompile("cn.com.imch.mddp.admin.aop.BatteryLogAspect-(.*)")
	matchArr := compileRegex.FindStringSubmatch(text)
	jsonData := matchArr[1]
	time := text[0:19]
	date := strings.Split(time, " ")[0]
	hour := strings.Split(time, " ")[1]
	var builder strings.Builder
	jsonData = strings.TrimRight(jsonData, "}")
	builder.WriteString(jsonData)
	builder.WriteString(`,"date":"`)
	builder.WriteString(date)
	builder.WriteString(`","time":"`)
	builder.WriteString(string(hour))
	builder.WriteString(`"}`)
	return builder.String()
}
