package examples

import (
	"log"
	"testing"
	
	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
	
	ffmpeg "github.com/ivanlebron/ffmpeg-go"
)

//
// More simple examples please refer to ffmpeg_test.go
//

func TestLive(t *testing.T) {
	
	if _, err := ffmpeg.Input("rtmp://localhost/live/test").
		Output("rtmp://localhost/live/mac",
			ffmpeg.KwArgs{
				"c:v": "h264",
				"c:a": "aac",
			},
			ffmpeg.KwArgs{
				"f": "flv",
			}).
		ErrorToStdOut().Run(func() {
		log.Println("1111")
		
	}); err != nil {
		log.Println("err ", err.Error())
	} else {
		log.Println("no err")
	}
	//assert.Nil(t, err)
	
	//ffmpeg.Input("").Output("").GlobalArgs().Run()
	
}

func TestExampleStream(t *testing.T) {
	ExampleStream("./sample_data/in1.mp4", "./sample_data/out1.mp4", false)
}

func TestExampleReadFrameAsJpeg(t *testing.T) {
	reader := ExampleReadFrameAsJpeg("./sample_data/in1.mp4", 5)
	img, err := imaging.Decode(reader)
	if err != nil {
		t.Fatal(err)
	}
	err = imaging.Save(img, "./sample_data/out1.jpeg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestExampleShowProgress(t *testing.T) {
	ExampleShowProgress("./sample_data/in1.mp4", "./sample_data/out2.mp4")
}

func TestExampleChangeCodec(t *testing.T) {
	_, err := ffmpeg.Input("./sample_data/in1.mp4").
		Output("./sample_data/out1.mp4", ffmpeg.KwArgs{"c:v": "libx265"}).
		OverWriteOutput().ErrorToStdOut().Run(nil)
	assert.Nil(t, err)
}

func TestExampleCutVideo(t *testing.T) {
	_, err := ffmpeg.Input("./sample_data/in1.mp4", ffmpeg.KwArgs{"ss": 1}).
		Output("./sample_data/out1.mp4", ffmpeg.KwArgs{"t": 1}).OverWriteOutput().Run(nil)
	assert.Nil(t, err)
}

func TestExampleScaleVideo(t *testing.T) {
	_, err := ffmpeg.Input("./sample_data/in1.mp4").
		Output("./sample_data/out1.mp4", ffmpeg.KwArgs{"vf": "scale=w=480:h=240"}).
		OverWriteOutput().ErrorToStdOut().Run(nil)
	assert.Nil(t, err)
}

func TestExampleAddWatermark(t *testing.T) {
	// show watermark with size 64:-1 in the top left corner after seconds 1
	overlay := ffmpeg.Input("./sample_data/overlay.png").Filter("scale", ffmpeg.Args{"64:-1"})
	_, err := ffmpeg.Filter(
		[]*ffmpeg.Stream{
			ffmpeg.Input("./sample_data/in1.mp4"),
			overlay,
		}, "overlay", ffmpeg.Args{"10:10"}, ffmpeg.KwArgs{"enable": "gte(t,1)"}).
		Output("./sample_data/out1.mp4").OverWriteOutput().ErrorToStdOut().Run(nil)
	assert.Nil(t, err)
}

func TestExampleCutVideoForGif(t *testing.T) {
	pid, err := ffmpeg.Input("./sample_data/in1.mp4", ffmpeg.KwArgs{"ss": "1"}).
		Output("./sample_data/out1.gif", ffmpeg.KwArgs{"s": "320x240", "pix_fmt": "rgb24", "t": "3", "r": "3"}).
		OverWriteOutput().ErrorToStdOut().Run(nil)
	log.Println("pid", pid)
	assert.Nil(t, err)
}

func TestExampleMultipleOutput(t *testing.T) {
	input := ffmpeg.Input("./sample_data/in1.mp4").Split()
	out1 := input.Get("0").Filter("scale", ffmpeg.Args{"1920:-1"}).
		Output("./sample_data/1920.mp4", ffmpeg.KwArgs{"b:v": "5000k"})
	out2 := input.Get("1").Filter("scale", ffmpeg.Args{"1280:-1"}).
		Output("./sample_data/1280.mp4", ffmpeg.KwArgs{"b:v": "2800k"})
	_, err := ffmpeg.MergeOutputs(out1, out2).OverWriteOutput().ErrorToStdOut().Run(nil)
	assert.Nil(t, err)
}
