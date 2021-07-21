local goodskey = tostring(KEYS[1])
print("goodskey", goodskey)
local num = tonumber(ARGV[1])
print("num", num)
local Fnum = 0 - num
local stock = redis.call("hget", goodskey, "stock")
if stock < num  then
    -- return 0
    print("here1")
    return 0;
end

redis.call("hincrby", goodskeys, "stock", Fnum)
-- return 1
print("here2")