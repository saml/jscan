local MultipartParser = require("resty.upload")
local cjson = require("cjson")

local function multipart_parse_key(header)
    _,_,key = string.find(header, 'name="([^"]+)"')
    return key
end

local function multipart_read(parser)
    local typ, res, err = parser:read()

    if not typ then
        mgx.say("failed to read: ", err)
        return false
    end

    if typ == 'eof' then
        return false
    end

    --if typ == 'header' then
    --    multipart_parse_key(res[1])
    
    ngx.say("read: ", typ, res)

    return true
end

local function multipart_to_json()
    --local parser = MultipartParser:new(4096)
    local parser = MultipartParser:new(5)
    parser:set_timeout(1000) -- 1 sec
    local is_normal_end = true
    while true do
        if not multipart_read(parser) then
            is_normal_end = false
            break
        end
    end

    if is_normal_end then
        --read remainder
        multipart_read(parser)
    end
end

local method = ngx.req.get_method()
if method == 'POST' then
    multipart_to_json()
    --ngx.req.read_body()
    --local args = ngx.req.get_post_args()
    --ngx.say(cjson.encode(args))
else
    ngx.status = ngx.HTTP_NOT_ALLOWED
    ngx.say(string.format('{"status": %d, "content": "Method Not Allowed: %s"}', ngx.HTTP_NOT_ALLOWED, method))
end
