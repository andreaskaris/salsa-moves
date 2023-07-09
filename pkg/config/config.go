package config

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"

	"gopkg.in/yaml.v3"
)

const (
	DefaultMoveList = "Mambo on 2"
)

type Moves struct {
	List map[string][]Move
}

type Move struct {
	Name   string
	Counts int
}

func (m Move) String() string {
	return fmt.Sprintf("%s (%d)", m.Name, m.Counts)
}

func ParseMove(m string) (Move, error) {
	re := regexp.MustCompile(`([a-zA-Z0-9/ ]+) \(([0-9]+)\)`)
	match := re.FindStringSubmatch(m)
	fmt.Println(match)
	if len(match) != 3 {
		return Move{}, fmt.Errorf("Invalid regex: %q, parsed to: %v", m, match)
	}
	counts, err := strconv.Atoi(match[2])
	if err != nil {
		return Move{}, err
	}
	return Move{Name: match[1], Counts: counts}, nil
}

type Song struct {
	BPM           int
	SleepForRand  int
	SleepForConst int
}

type Text struct {
	Size int
}

type Config struct {
	Moves                     Moves
	MaxMoves                  int
	MinMoves                  int
	ProbabilityRandomSequence int
	Text                      Text
	Song                      Song
	fileName                  string
}

func Load(fileName string) (*Config, error) {
	c := &Config{}
	fContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s, err: %q", fileName, err)
	}
	if err := yaml.Unmarshal(fContent, c); err != nil {
		return nil, fmt.Errorf("error parsing YAML file %s, err: %q", fileName, err)
	}
	c.fileName = fileName
	return c, nil
}

func (c *Config) Save() error {
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	fmt.Println(c.fileName, string(out))
	return os.WriteFile(c.fileName, out, 0640)
}

// AddMove adds the provided move to the provided list.
func (c *Config) AddMove(listName string, move Move) {
	if c.Moves.List == nil {
		c.Moves.List = make(map[string][]Move)
	}
	c.Moves.List[listName] = append(c.Moves.List[listName], move)
	sort.Slice(c.Moves.List[listName], func(i, j int) bool {
		return c.Moves.List[listName][i].Name < c.Moves.List[listName][j].Name
	})
}

// DeleteMove deletes the provided move from the provided list.
func (c *Config) DeleteMove(listName string, name string) {
	for i := range c.Moves.List[listName] {
		if c.Moves.List[listName][i].Name == name {
			c.Moves.List[listName] = append(c.Moves.List[listName][:i], c.Moves.List[listName][i+1:]...)
			return
		}
	}
}

func (c *Config) GetMoveList(listName string) []Move {
	return c.Moves.List[listName]
}

func (c *Config) GetMoveStringList(listName string) []string {
	var moves []string
	for _, m := range c.GetMoveList(listName) {
		moves = append(moves, m.String())
	}
	return moves
}

func (c *Config) SetMaxMoves(max int) {
	c.MaxMoves = max
}

func (c *Config) GetMaxMoves() int {
	return c.MaxMoves
}

func (c *Config) SetMinMoves(min int) {
	c.MinMoves = min
}

func (c *Config) GetMinMoves() int {
	return c.MinMoves
}

func (c *Config) SetProbabilityRandomSequence(p int) {
	c.ProbabilityRandomSequence = p
}

func (c *Config) GetProbabilityRandomSequence() int {
	return c.ProbabilityRandomSequence
}

func (c *Config) GetBPM() int {
	return c.Song.BPM
}

func (c *Config) SetBPM(bpm int) {
	c.Song.BPM = bpm
}

func (c *Config) GetSleepForRand() int {
	return c.Song.SleepForRand
}

func (c *Config) SetSleepForRand(sleep int) {
	c.Song.SleepForRand = sleep
}

func (c *Config) GetSleepForConst() int {
	return c.Song.SleepForConst
}

func (c *Config) SetSleepForConst(sleep int) {
	c.Song.SleepForConst = sleep
}

func (c *Config) GetTextSize() int {
	return c.Text.Size
}

func (c *Config) SetTextSize(size int) {
	c.Text.Size = size
}
