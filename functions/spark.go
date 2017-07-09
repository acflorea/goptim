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

	dataset := "netbeans"

	sparkSubmit := "/Users/acflorea/Bin/spark-1.6.2-bin-hadoop2.6/bin/spark-submit"
	targetJar := "/Users/acflorea/phd/columbugus/target/scala-2.10/columbugus-assembly-2.3.1.jar"

	configFile := "/Users/acflorea/Bin/spark-1.6.2-bin-hadoop2.6/columbugus-conf/" + dataset + ".conf"
	fsRoot := "/Users/acflorea/phd/columbugus_data/" + dataset + "_final_test"

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
		p.Values["categoryScalingFactor"].(string) +
		" -Dreccsys.preprocess.productScalingFactor=" +
		p.Values["productScalingFactor"].(string) +
		" -Dreccsys.train.stepSize=" +
		"1" +
		" -Dreccsys.train.regParam=" +
		"0.01"

	//"--master=local[*]"
	//"--executor-memory 10G"
	//"--driver-memory 64G"
	//"--driver-java-options -Xmx80G"
	//"--driver-java-options" + sparkParams

	//cmd := exec.Command(sparkSubmit, sparkParams, "&2 > ./xxx.log")
	cmd := exec.Command(sparkSubmit,
		"--class", "dr.acf.recc.ReccomenderBackbone",
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
	fmt.Println(p.Values["categoryScalingFactor"].(string), p.Values["productScalingFactor"].(string), string(dat))

	f1str := strings.TrimPrefix(strings.Split(resultsStr, " ")[2], "F:")
	f1Measure, _ := strconv.ParseFloat(f1str, 64)

	return f1Measure, nil
}
