package model

type Node struct {
	ID         int `json:"id"`
	Label      string
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	Sex        int64  `json:"sex"`
	City       string `json:"city"`
}

type Relation struct {
	Node    Node
	RelType string
	EndNode Node
}

type RelationDTO struct {
	RelType   string `json:"rel_type"`
	EndNodeID int    `json:"end_node_id"`
}

type InsertRequest struct {
	Node      Node          `json:"node"`
	Relations []RelationDTO `json:"relations"`
}
