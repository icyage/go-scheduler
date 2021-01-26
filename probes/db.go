package probes

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/obsurvive/voyager/utils"
	"time"
)

type Probe struct {
	UUID string
	Url  string
}

type Result struct {
	createdAt     string
	updatedAt     string
	statusCode    int
	error         string
	content       string
	latency       int
	sslExpireDate string
	brokenLinks   string
	uuid          string
}

func GetProbe(uuid string) (Probe, error) {
	cache := utils.GetCache()
	cachedProbe, err := cache.Get(uuid).Result()
	if err == redis.Nil {
		err, db := utils.GetDB()
		if err != nil {
			return Probe{}, nil
		}
		// grab from DB
		row := db.QueryRow("SELECT uuid, url FROM public.probe WHERE uuid='982e50f7-48cc-40b0-95a3-bbf8b53ce229'")

		probe := Probe{}
		err = row.Scan(&probe.UUID, &probe.Url)
		if err != nil {
			return Probe{}, err
		}

		// pop cache
		probeJson, _ := json.Marshal(probe)
		err = cache.Set(probe.UUID, string(probeJson), cfg.GetDuration("redis_expire_seconds") * time.Second ).Err()
		if err != nil {
			return Probe{}, err
		}
		// return probe
		return probe, nil
	} else if err != nil {
		return Probe{}, err
	} else {
		probe := Probe{}
		if err := json.Unmarshal([]byte(cachedProbe), &probe); err != nil {
			return Probe{}, err
		}
		return probe, nil
	}
}
