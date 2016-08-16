package endpoint

import (
	"fmt"
	"log"
	"strings"
)

func (e *endpoint) sendRegistrationRequest(action string) (string, error) {
	request := "<?xml version=\"1.0\" encoding=\"UTF-8\" ?>"
	request += "<root>"
	request += fmt.Sprintf(
		"<S2PC_Regi UserID=\"%s\" UniqueID=\"%s\" RegiType="+
			"\"%s\" />", e.name, e.uid, action)
	request += "</root>"

	return formUpload(e.s.IP(), "/IDS/ScanFaxToPC.cgi", "c:\\IDS.XML",
		request, false)
}

func parseRegistrationResponse(res string) (string, error) {
	data := strings.Split(res, "\"")

	result, instance := "", ""
	for i := 0; i < len(data); i++ {
		if data[i] == " Result=" {
			result = strings.TrimSpace(data[i+1])
		} else if data[i] == " InstanceID=" {
			instance = strings.TrimSpace(data[i+1])
		}
	}

	if (result != "ADD_OK" && result != "DELETE_OK") || len(instance) == 0 {
		return "", fmt.Errorf("error parsing response: %q", result)
	}

	return instance, nil
}

func (e *endpoint) register() error {
	e.m.Lock()
	defer e.m.Unlock()

	log.Printf("%s %q Registering endpoint.", e.uid, e.name)

	res, err := e.sendRegistrationRequest("ADD")
	if err != nil {
		log.Println("Error sending registration request.")
		return err
	}

	log.Println("Sent registration request. No error.")

	instance, err := parseRegistrationResponse(res)
	if err != nil {
		return err
	}

	e.instance = instance

	log.Printf("%s %q Instance ID = %s.", e.uid, e.name, e.instance)

	return nil
}

func (e *endpoint) unregister() error {
	e.m.Lock()
	defer e.m.Unlock()

	log.Printf("%s Unregistering endpoint.", e.uid)

	res, err := e.sendRegistrationRequest("DELETE")
	if err != nil {
		return err
	}

	instance, err := parseRegistrationResponse(res)
	if err != nil {
		return err
	}

	if instance != e.instance {
		return fmt.Errorf(
			"unexpected instance id. Got %s, expected %s.",
			instance, e.instance)
	}

	e.instance = ""

	return nil
}
