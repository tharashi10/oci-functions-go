package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"func/stream"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	fdk "github.com/fnproject/fdk-go"
	"github.com/oracle/oci-go-sdk/v55/common/auth"
	"github.com/oracle/oci-go-sdk/v55/example/helpers"
	"github.com/oracle/oci-go-sdk/v55/objectstorage"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(targzHandler))
}

func targzHandler(ctx context.Context, in io.Reader, out io.Writer) {
	// Get Message from Stream
	var msg []stream.Message
	var events stream.Events

	str := StreamToString(in)
	json.Unmarshal([]byte(str), &msg)

	log.Println("[INFO] Get ", len(msg), " Messages.")
	log.Print("===================================================")
	log.Println("[INFO] Get Msg Array ", msg)
	log.Println("[INFO] Get 1st Element", msg[0])
	log.Println("message: [Base64]", msg[0].Value)

	decoded, err := base64.StdEncoding.DecodeString(msg[0].Value)
	if err != nil {
		log.Fatal(err)
	}

	// Decode values
	fmt.Println("message: [Base64_decoded]", string(decoded))
	json.Unmarshal([]byte(string(decoded)), &events)

	// Debug
	log.Println("events: [Namespace]", events.Data.AdditionalDetails.Namespace)
	log.Println("message: [Bucket]", events.Data.AdditionalDetails.BucketName)
	log.Println("message: [ResourceName]", events.Data.ResourceName)
	log.Print("===================================================")

	// Parameter
	namespace := events.Data.AdditionalDetails.Namespace
	bucketName := events.Data.AdditionalDetails.BucketName
	resourceName := events.Data.ResourceName

	trgNamespace := events.Data.AdditionalDetails.Namespace
	trgBucketName := events.Data.AdditionalDetails.BucketName
	trgPrefix := "temp-obj"

	log.Println("events: [trgNamespace]", trgNamespace)
	log.Println("message: [trgBucket]", trgBucketName)
	log.Println("message: [trgResourceName]", trgPrefix)
	log.Print("===================================================")

	// Prepare Resource Principal
	provider, err := auth.ResourcePrincipalConfigurationProvider()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create Gor and Get tar object
	objectStorageClient, clerr := objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
	helpers.FatalIfError(clerr)
	resp := getTarObject(ctx, objectStorageClient, namespace, bucketName, resourceName)
	content, _ := ioutil.ReadAll(resp.Content)
	fmt.Println("Bytes: ", content)

	r := bytes.NewReader(content)
	fmt.Println("StringFromBytes: ", StreamToString(r))

	// Input : sample.tar.gz, io.reader
	// Tar解凍開始
	gzr, err := gzip.NewReader(r)
	if err != nil {
		fmt.Println("StringFromBytes: ", StreamToString(r))
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	// Headerごとに中身解凍(Open and iterate through the files in the archive.)
	for {
		header, err := tr.Next()
		fmt.Println("HeaderName : ", header)
		switch {
		case err == io.EOF:
			fmt.Println(err)
			//return nil

		case err != nil:
			fmt.Println(err)
			//return err

		case header == nil:
			continue
		}
		fmt.Printf("Contents of %s:\n", header.Name)

		//resourceName作成
		target := filepath.Join(trgPrefix, header.Name)
		fmt.Println("HeaderName :", header.Name) //HeaderName : sample.csv
		fmt.Println("HeaderNameFullpath :", target)
		fmt.Println("====================")

		bs, _ := ioutil.ReadAll(tr)
		// convert the []byte to a string
		s := string(bs)
		fmt.Println("Bytes:", bs)
		fmt.Println("String:", s)

		// Passing Object Storage
		putObject(ctx, objectStorageClient, namespace, bucketName, target, s, nil)
		log.Println("Task has Done.")
	}
}

func getTarObject(ctx context.Context, client objectstorage.ObjectStorageClient, namespace, bucketName, resourceName string) objectstorage.GetObjectResponse {
	req := objectstorage.GetObjectRequest{NamespaceName: &namespace,
		BucketName: &bucketName,
		ObjectName: &resourceName}
	resp, err := client.GetObject(context.Background(), req)
	helpers.FatalIfError(err)
	fmt.Println(resp)
	return resp
}

func putObject(ctx context.Context, c objectstorage.ObjectStorageClient, namespace, bucketname, objectname string, objectContent string, metadata map[string]string) error {
	helper := int64(len(objectContent))
	type Str struct {
		v *int64
	}
	inst := Str{
		v: &helper,
	}
	request := objectstorage.PutObjectRequest{
		NamespaceName: &namespace,
		BucketName:    &bucketname,
		ObjectName:    &objectname,
		ContentLength: inst.v,
		PutObjectBody: ioutil.NopCloser(strings.NewReader(objectContent)),
		OpcMeta:       metadata,
	}
	_, err := c.PutObject(ctx, request)
	fmt.Println("put object")
	return err
}

func getOciNamespace(ctx context.Context, client objectstorage.ObjectStorageClient) string {
	request := objectstorage.GetNamespaceRequest{}
	r, err := client.GetNamespace(ctx, request)
	helpers.FatalIfError(err)
	fmt.Println("get namespace Done.")
	return *r.Value
}

func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String()
}
