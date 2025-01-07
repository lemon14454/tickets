package cache

import "github.com/redis/go-redis/v9"

func NewRedis(address string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}

var ClaimCacheTicket = redis.NewScript(`
local ticket_ids = ARGV
local quantity = tonumber(KEYS[1])
local userID = tonumber(KEYS[2])

local claimed_tickets = {}
local count = 0

for i, ticket_id in ipairs(ticket_ids) do
    local key = "ticket:" .. ticket_id

    local result = redis.call("SETNX", key, userID)
    if result == 1 then
        redis.call("EXPIRE", key, 300)

        table.insert(claimed_tickets, ticket_id)
        count = count + 1

        if count == quantity then
            break
        end
    end
end

if count < quantity then
    for _, ticket_id in ipairs(claimed_tickets) do
        local key = "ticket:" .. ticket_id
        redis.call("DEL", key)
    end
    return {}
end

return claimed_tickets
`)
