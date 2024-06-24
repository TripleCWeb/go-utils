package leaderboard

import (
	"context"
	"fmt"
	"math"

	"github.com/redis/go-redis/v9"
)

type User struct {
	Name  string
	Score int64
	Rank  int64
}

type Team struct {
	Name    string
	Members map[string]User
	Rank    int
}

type RedisSettings struct {
	Host     string
	Password string
}

type Leaderboard struct {
	redis    *redis.Client
	Name     string
	PageSize int
}

func NewLeaderboard(redis *redis.Client, name string, pageSize int) Leaderboard {
	l := Leaderboard{redis: redis, Name: name, PageSize: pageSize}
	return l
}

func (l *Leaderboard) Clear() error {
	err := l.redis.Del(context.Background(), l.Name).Err()
	if err != nil {
		fmt.Printf("error on destroy Leaderboard:%s", l.Name)
	}
	return err
}

func (l *Leaderboard) RankMember(username string, score int64) (u *User, err error) {
	if l.redis.ZAdd(context.Background(), l.Name, redis.Z{Score: float64(score), Member: username}).Err() != nil {
		fmt.Printf("error on store in redis in rankMember Leaderboard:%s - Username:%s - Score:%d", l.Name, username, score)
		return
	}
	rank := l.GetRank(username)
	u = &User{Name: username, Score: score, Rank: rank}
	return
}

func (l *Leaderboard) RankMemberInc(username string, score int64) (u *User, err error) {
	if l.redis.ZIncrBy(context.Background(), l.Name, float64(score), username).Err() != nil {
		fmt.Printf("error on store in redis in rankMember Leaderboard:%s - Username:%s - Score:%d", l.Name, username, score)
		return
	}

	rank := l.GetRank(username)
	u = &User{Name: username, Score: score, Rank: rank}
	return
}

func (l *Leaderboard) TotalMembers() int64 {
	res := l.redis.ZCard(context.Background(), l.Name)
	if res.Err() != nil {
		fmt.Printf("error on get leaderboard total members")
		return 0
	}
	return res.Val()
}

func (l *Leaderboard) RemoveMember(username string) (u *User, err error) {
	u, err = l.GetMember(username)
	res := l.redis.ZRem(context.Background(), l.Name, username)
	if res.Err() != nil {
		fmt.Printf("error on remove user from leaderboard")
	}
	return
}

func (l *Leaderboard) TotalPages() int {
	pages := 0

	var total int64
	res := l.redis.ZCount(context.Background(), l.Name, "-inf", "+inf")
	if res.Err() == nil {
		total = res.Val()
		pages = int(math.Ceil(float64(total) / float64(l.PageSize)))
	}
	return pages
}

func (l *Leaderboard) GetMember(username string) (u *User, err error) {
	var score int64
	{
		res := l.redis.ZScore(context.Background(), l.Name, username)
		err = res.Err()
		if err != nil {
			if err == redis.Nil {
				err = nil
				return
			} else {
				fmt.Printf("error on get user rank Leaderboard:%s - Username:%s", l.Name, username)
				return
			}
		}
		score = int64(res.Val())
	}

	rank := l.GetRank(username)
	u = &User{Name: username, Score: score, Rank: rank}
	return
}

func (l *Leaderboard) GetAroundMe(username string) []*User {
	currentUser, _ := l.GetMember(username)
	startOffset := int(currentUser.Rank) - (l.PageSize / 2)
	if startOffset < 0 {
		startOffset = 0
	}
	endOffset := (startOffset + l.PageSize) - 1
	return l.GetMembersByRange(startOffset, endOffset)
}

func (l *Leaderboard) GetMembersByRange(startOffset int, endOffset int) (users []*User) {
	res := l.redis.ZRevRangeWithScores(context.Background(), l.Name, int64(startOffset), int64(endOffset))

	for _, value := range res.Val() {
		username := value.Member.(string)
		rank := l.GetRank(username)
		users = append(users, &User{Name: username, Score: int64(value.Score), Rank: rank})
	}
	return users
}

func (l *Leaderboard) GetRank(username string) int64 {
	res := l.redis.ZRevRank(context.Background(), l.Name, username)
	if res.Err() != nil {
		fmt.Printf("Leaderboard::GetRank error on get user rank Leaderboard:%s - Username:%s\n", l.Name, username)
		return -1
	}
	return res.Val() + 1
}

func (l *Leaderboard) GetLeaders(page int) []*User {
	if page < 1 {
		page = 1
	}
	if page > l.TotalPages() {
		page = l.TotalPages()
	}
	redisIndex := page - 1
	startOffset := redisIndex * l.PageSize
	if startOffset < 0 {
		startOffset = 0
	}
	endOffset := (startOffset + l.PageSize) - 1
	return l.GetMembersByRange(startOffset, endOffset)
}

func (l *Leaderboard) GetMemberByRank(position int64) (u *User) {
	if position <= l.TotalMembers() {
		currentPage := int(math.Ceil(float64(position) / float64(l.PageSize)))
		offset := (position - 1) % int64(l.PageSize)
		leaders := l.GetLeaders(currentPage)
		if leaders[offset].Rank == position {
			return leaders[offset]
		}
	}
	return
}

/* End Public functions */
