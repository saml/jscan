
local method = ngx.req.get_method()
if method == 'POST' then
    ngx.req.read_body()
    local args = ngx.req.get_post_args()
    ngx.say(cjson.encode(args))
else
    ngx.status = ngx.HTTP_NOT_ALLOWED
    ngx.say(string.format('{"status": %d, "content": "Method Not Allowed: %s"}', ngx.HTTP_NOT_ALLOWED, method))
end
