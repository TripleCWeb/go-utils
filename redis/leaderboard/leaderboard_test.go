package leaderboard

import (
	"strconv"
	"testing"

	"github.com/redis/go-redis/v9"

	"launchpad.net/gocheck"
)

func Test(t *testing.T) {
	gocheck.TestingT(t)
}

type S struct{}

var _ = gocheck.Suite(&S{})

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "host.docker.internal:6379",
	Password: "rdspwd11131456",
	DB:       int(4),
})

func (s *S) TestRankMember(c *gocheck.C) {
	highScore := NewLeaderboard(redisClient, "highScore", 10)
	dayvson, err := highScore.RankMember("dayvson", 481516)
	c.Assert(err, gocheck.IsNil)
	arthur, err := highScore.RankMember("arthur", 1000)
	c.Assert(err, gocheck.IsNil)

	c.Assert(dayvson.Rank, gocheck.Equals, int64(1))
	c.Assert(arthur.Rank, gocheck.Equals, int64(2))

	c.Assert(highScore.GetRank(dayvson.Name), gocheck.Equals, int64(1))
	c.Assert(highScore.GetRank(arthur.Name), gocheck.Equals, int64(2))
	highScore.Clear()
}

func (s *S) TestRankMemberInc(c *gocheck.C) {
	// rank
	highScore := NewLeaderboard(redisClient, "highScoreInc", 10)
	dayvson, err := highScore.RankMemberInc("dayvson", 481516)
	c.Assert(err, gocheck.IsNil)
	arthur, err := highScore.RankMemberInc("arthur", 1000)
	c.Assert(err, gocheck.IsNil)
	c.Assert(dayvson.Rank, gocheck.Equals, int64(1))
	c.Assert(arthur.Rank, gocheck.Equals, int64(2))

	// add
	arthur, err = highScore.RankMemberInc("arthur", 1000000)
	c.Assert(err, gocheck.IsNil)

	// change rank
	c.Assert(arthur.Rank, gocheck.Equals, int64(1))

	highScore.Clear()
}

func (s *S) TestTotalMembers(c *gocheck.C) {
	bestTime := NewLeaderboard(redisClient, "bestTime", 10)
	for i := 0; i < 10; i++ {
		bestTime.RankMember("member_"+strconv.Itoa(i), int64(1234*i))
	}
	totalMembers := bestTime.TotalMembers()
	c.Assert(totalMembers, gocheck.Equals, int64(10))
	bestTime.Clear()
}

func (s *S) TestRemoveMember(c *gocheck.C) {
	bestTime := NewLeaderboard(redisClient, "bestWeek", 10)
	for i := 0; i < 10; i++ {
		bestTime.RankMember("member_"+strconv.Itoa(i), int64(1234*i))
	}
	c.Assert(bestTime.TotalMembers(), gocheck.Equals, int64(10))
	bestTime.RemoveMember("member_5")
	c.Assert(bestTime.TotalMembers(), gocheck.Equals, int64(9))
	bestTime.Clear()
}

func (s *S) TestTotalPages(c *gocheck.C) {
	bestTime := NewLeaderboard(redisClient, "All", 25)
	for i := 0; i < 101; i++ {
		bestTime.RankMember("member_"+strconv.Itoa(i), int64(1234*i))
	}
	c.Assert(bestTime.TotalPages(), gocheck.Equals, 5)
	bestTime.Clear()
}

func (s *S) TestGetUser(c *gocheck.C) {
	friendScore := NewLeaderboard(redisClient, "friendScore", 10)
	dayvson, _ := friendScore.RankMember("dayvson", 12345)
	felipe, _ := friendScore.RankMember("felipe", 12344)
	c.Assert(dayvson.Rank, gocheck.Equals, int64(1))
	c.Assert(felipe.Rank, gocheck.Equals, int64(2))
	friendScore.RankMember("felipe", 12346)
	felipe, _ = friendScore.GetMember("felipe")
	dayvson, _ = friendScore.GetMember("dayvson")
	c.Assert(felipe.Rank, gocheck.Equals, int64(1))
	c.Assert(dayvson.Rank, gocheck.Equals, int64(2))
	friendScore.Clear()
}

func (s *S) TestGetAroundMe(c *gocheck.C) {
	bestTime := NewLeaderboard(redisClient, "BestAllTime", 25)
	for i := 0; i < 101; i++ {
		bestTime.RankMember("member_"+strconv.Itoa(i), int64(1234*i))
	}
	users := bestTime.GetAroundMe("member_20")
	firstAroundMe := users[0]
	lastAroundMe := users[bestTime.PageSize-1]
	c.Assert(len(users), gocheck.Equals, bestTime.PageSize)
	c.Assert(firstAroundMe.Name, gocheck.Equals, "member_31")
	c.Assert(lastAroundMe.Name, gocheck.Equals, "member_7")
	bestTime.Clear()
}

func (s *S) TestGetRank(c *gocheck.C) {
	sevenDays := NewLeaderboard(redisClient, "7days", 25)
	for i := 0; i < 101; i++ {
		sevenDays.RankMember("member_"+strconv.Itoa(i), int64(1234*i))
	}
	sevenDays.RankMember("member_6", 1000)
	c.Assert(sevenDays.GetRank("member_6"), gocheck.Equals, int64(100))
	sevenDays.Clear()
}

func (s *S) TestGetLeaders(c *gocheck.C) {
	bestYear := NewLeaderboard(redisClient, "bestYear", 25)
	for i := 0; i < 1000; i++ {
		bestYear.RankMember("member_"+strconv.Itoa(i+1), int64(1234*i))
	}
	var users = bestYear.GetLeaders(1)

	firstOnPage := users[0]
	lastOnPage := users[len(users)-1]
	c.Assert(len(users), gocheck.Equals, bestYear.PageSize)
	c.Assert(firstOnPage.Name, gocheck.Equals, "member_1000")
	c.Assert(firstOnPage.Rank, gocheck.Equals, int64(1))
	c.Assert(lastOnPage.Name, gocheck.Equals, "member_976")
	c.Assert(lastOnPage.Rank, gocheck.Equals, int64(25))
	bestYear.Clear()
}

func (s *S) TestGetUserByRank(c *gocheck.C) {
	sevenDays := NewLeaderboard(redisClient, "week", 25)
	for i := 0; i < 101; i++ {
		sevenDays.RankMember("member_"+strconv.Itoa(i), int64(1234*i))
	}
	member := sevenDays.GetMemberByRank(10)
	c.Assert(member.Name, gocheck.Equals, "member_91")
	c.Assert(member.Rank, gocheck.Equals, int64(10))
	sevenDays.Clear()
}
