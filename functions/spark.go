package functions

import (
	"os/exec"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
)

// A more complicated function (submits a task to Apache Spark, and gathers the results)
func SparkIt(p MultidimensionalPoint, vargs map[string]interface{}) (float64, error) {

	mainClass := vargs["mainClass"].(string) // "dr.acf.recc.ReccomenderBackbone"

	sparkSubmit := vargs["sparkSubmit"].(string) // "/Users/acflorea/Bin/spark-1.6.2-bin-hadoop2.6/bin/spark-submit"
	targetJar := vargs["targetJar"].(string)     // "/Users/acflorea/phd/columbugus/target/scala-2.10/columbugus-assembly-2.3.1.jar"

	configFile := vargs["configFile"].(string) // "/Users/acflorea/Bin/spark-1.6.2-bin-hadoop2.6/columbugus-conf/netbeans.conf"
	fsRoot := vargs["fsRoot"].(string)         //"/Users/acflorea/phd/columbugus_data/netbeans_final_test"

	resultsFileName := "gorand_results.out"
	tuningMode := "true"

	sparkParams := "-Dconfig.file=" +
		configFile +
		" " +
		"-Dreccsys.phases.preprocess=true " +
		"-Dreccsys.preprocess.includeCategory=true " +
		"-Dreccsys.preprocess.includeProduct=true " +
		"-Dreccsys.global.tuningMode=" + tuningMode +
		" -Dreccsys.filesystem.resultsFileName=" +
		resultsFileName +
		" -Dreccsys.filesystem.root=" +
		fsRoot +
		" -Dreccsys.preprocess.categoryScalingFactor=" +
		FloatToString(p.Values["categoryScalingFactor"].(float64)) +
		" -Dreccsys.preprocess.productScalingFactor=" +
		FloatToString(p.Values["productScalingFactor"].(float64)) +
		" -Dreccsys.train.stepSize=" +
		"1" +
		" -Dreccsys.train.regParam=" +
		FloatToString(p.Values["C"].(float64))

	//"--master=local[*]"
	//"--executor-memory 10G"
	//"--driver-memory 64G"
	//"--driver-java-options -Xmx80G"
	//"--driver-java-options" + sparkParams

	//cmd := exec.Command(sparkSubmit, sparkParams, "&2 > ./xxx.log")
	cmd := exec.Command(sparkSubmit,
		"--class", mainClass,
		"--master=local[3]",
		"--driver-java-options",
		sparkParams,
		targetJar)

	if err := cmd.Run(); err != nil {
		fmt.Println("Error: ", err)
		return 0.0, err
	}

	dat, _ := ioutil.ReadFile(fsRoot + "/" + resultsFileName)
	resultsStr := string(dat)
	fmt.Println(
		FloatToString(p.Values["C"].(float64)),
		FloatToString(p.Values["categoryScalingFactor"].(float64)),
		FloatToString(p.Values["productScalingFactor"].(float64)),
		string(dat))

	f1str := strings.TrimPrefix(strings.Split(resultsStr, " ")[2], "F:")
	f1Measure, _ := strconv.ParseFloat(f1str, 64)

	return f1Measure, nil
}
