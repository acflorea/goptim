package functions

import (
	"os/exec"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
)

// A more complicated function (submits a task to Apache Spark, and gathers the results)
func SparkIt(p MultidimensionalPoint) (float64, error) {

	sparkSubmit := "/Users/acflorea/Bin/spark-1.6.2-bin-hadoop2.6/bin/spark-submit"
	configFile := "/Users/acflorea/Bin/spark-1.6.2-bin-hadoop2.6/columbugus-conf/netbeans.conf"
	fsRoot := "/Users/acflorea/phd/columbugus_data/netbeans_final_test"
	targetJar := "/Users/acflorea/phd/columbugus/target/scala-2.10/columbugus-assembly-2.1.jar"
	resultsFileName := "gorand_results.out"
	mode := "true"

	sparkParams := "-Dconfig.file=" +
		configFile +
		" " +
		"-Dreccsys.phases.preprocess=true " +
		"-Dreccsys.preprocess.includeCategory=true " +
		"-Dreccsys.preprocess.includeProduct=true " +
		"-Dreccsys.global.tuningMode=" + mode +
		" -Dreccsys.filesystem.resultsFileName=" +
		resultsFileName +
		" -Dreccsys.filesystem.root=" +
		fsRoot +
		" -Dreccsys.preprocess.categoryScalingFactor=" +
		FloatToString(p.Values[0]) +
		" -Dreccsys.preprocess.productScalingFactor=" +
		FloatToString(p.Values[1]) +
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
	cmd := exec.Command(sparkSubmit, "--class", "dr.acf.recc.ReccomenderBackbone", "--driver-java-options", sparkParams, targetJar)

	if err := cmd.Run(); err != nil {
		fmt.Println("Error: ", err)
		return 0.0, err
	}

	dat, _ := ioutil.ReadFile(fsRoot + "/" + resultsFileName)
	resultsStr := string(dat)
	fmt.Println(FloatToString(p.Values[0]), FloatToString(p.Values[1]), string(dat))

	f1str := strings.TrimPrefix(strings.Split(resultsStr, " ")[2], "F:")
	f1Measure, _ := strconv.ParseFloat(f1str, 64)

	return f1Measure, nil
}
