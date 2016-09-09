package cascadestore

import (
    "encoding/gob"
    "encoding/json"

    "fmt"
    "bytes"

    "github.com/gorilla/sessions"
)

// SessionSerializer provides an interface hook for alternative serializers
type SessionSerializer interface {
    Deserialize(d []byte, ss *sessions.Session) error
    Serialize(ss *sessions.Session) ([]byte, error)
}

// JSONSerializer encode the session map to JSON.
type JSONSerializer struct{}

// Serialize to JSON. Will err if there are unmarshalable key values
func (s JSONSerializer) Serialize(ss *sessions.Session) ([]byte, error) {
    m := make(map[string]interface{}, len(ss.Values))
    for k, v := range ss.Values {
        ks, ok := k.(string)
        if !ok {
            err := fmt.Errorf("Non-string key value, cannot serialize session to JSON: %v", k)
            return nil, err
        }
        m[ks] = v
    }
    return json.Marshal(m)
}

// Deserialize back to map[string]interface{}
func (s JSONSerializer) Deserialize(d []byte, ss *sessions.Session) error {
    m := make(map[string]interface{})
    err := json.Unmarshal(d, &m)
    if err != nil {
        fmt.Printf("redistore.JSONSerializer.deserialize() Error: %v", err)
        return err
    }
    for k, v := range m {
        ss.Values[k] = v
    }
    return nil
}

// GobSerializer uses gob package to encode the session map
type GobSerializer struct{}

// Serialize using gob
func (s GobSerializer) Serialize(ss *sessions.Session) ([]byte, error) {
    buf := new(bytes.Buffer)
    enc := gob.NewEncoder(buf)
    err := enc.Encode(ss.Values)
    if err == nil {
        return buf.Bytes(), nil
    }
    return nil, err
}

// Deserialize back to map[interface{}]interface{}
func (s GobSerializer) Deserialize(d []byte, ss *sessions.Session) error {
    dec := gob.NewDecoder(bytes.NewBuffer(d))
    return dec.Decode(&ss.Values)
}
