package config

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Moves struct {
	List []Move
	Max  int
	Min  int
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
	Moves    Moves
	Text     Text
	Song     Song
	fileName string
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

func (c *Config) AddMove(move Move) {
	if c.Moves.List == nil {
		c.Moves.List = []Move{}
	}
	c.Moves.List = append(c.Moves.List, move)
	sort.Slice(c.Moves.List, func(i, j int) bool {
		return c.Moves.List[i].Name < c.Moves.List[j].Name
	})
}

func (c *Config) DeleteMove(name string) {
	moves := c.GetMoveList()
	for i := range moves {
		if moves[i].Name == name {
			c.Moves.List = append(moves[:i], moves[i+1:]...)
			return
		}
	}
}

func (c *Config) GetMoveList() []Move {
	return c.Moves.List
}

func (c *Config) GetMoveStringList() []string {
	var moves []string
	for _, m := range c.GetMoveList() {
		moves = append(moves, m.String())
	}
	return moves
}

func (c *Config) SetMaxMoves(max int) {
	c.Moves.Max = max
}

func (c *Config) GetMaxMoves() int {
	return c.Moves.Max
}

func (c *Config) SetMinMoves(min int) {
	c.Moves.Min = min
}

func (c *Config) GetMinMoves() int {
	return c.Moves.Min
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
