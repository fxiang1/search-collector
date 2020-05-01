
package transforms

import (
	"strings"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)


type InsightResource struct {
	*unstructured.Unstructured
}

func (i InsightResource) BuildNode() Node {
	node := transformCommon(i)         // Start off with the common properties

	// apiGroupVersion(i.TypeMeta, &node) // add kind, apigroup and version
	node.Properties["kind"] = i.UnstructuredContent()["kind"]
	apiVersion := strings.Split(i.UnstructuredContent()["apiVersion"].(string), "/")
	if len(apiVersion) == 2 {
		node.Properties["apigroup"] = apiVersion[0]
		node.Properties["apiversion"] = apiVersion[1]
	} else {
		node.Properties["apiversion"] = apiVersion[0]
	}

	// Extract the properties specific to this type
	spec, _, _ := unstructured.NestedMap(i.UnstructuredContent(), "spec")
	problem, _, _ := unstructured.NestedMap(spec, "problem")
	description, _, _ := unstructured.NestedString(problem, "description")
	
	solutions, _, _ := unstructured.NestedSlice(spec, "solutions")
	topSolution := solutions[0].(map[string]interface{})

	resolution,_,_ := unstructured.NestedString(spec, "resolution")
	

	node.Properties["description"] = description
	node.Properties["topsolution"] = topSolution["description"]
	node.Properties["confidence"] = topSolution["confidence"]
	node.Properties["resolution"] = resolution

	return node
}

func (i InsightResource) BuildEdges(ns NodeStore) []Edge {
	//no op for now to implement interface
	return []Edge{}
}
