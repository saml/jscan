local cjson = require("cjson")
local RestyUpload = require('resty.upload')
local say = ngx.say
local request = ngx.req
local method = request.get_method()

local LINE_SEPARATOR = ' ' -- U+2028
local PARAGRAPH_SEPARATOR = ' ' -- U+2029
        
local extract_from_multipart = function(target, chunk_size, timeout)
    chunk_size = chunk_size or 4096
    timeout = timeout or 2000
    
    local body = RestyUpload:new(chunk_size)
    body:set_timeout(timeout)
    local collect = nil
    while true do
        local kind, data, err = body:read()
        if not kind then
            break
        end

        if kind == 'header' then
            local what = data[1]:lower()
            if what == 'content-disposition' and data[2]:find(target) then
                collect = {} -- start collecting
            end
        elseif kind == 'body' and collect then
            table.insert(collect, data)--:gsub(LINE_SEPARATOR, '\\u2028'):gsub(PARAGRAPH_SEPARATOR, '\\u2029'))
        elseif kind == 'part_end' and collect then
            break
        elseif kind == 'eof' then
            break
        end

    end
    if collect then
        return table.concat(collect)
    end
end


if method == 'POST' then
    local data = extract_from_multipart('name="json"')
    local obj = cjson.decode(data)
    local json = cjson.encode(obj)
    say(json)
else
    ngx.status = ngx.HTTP_NOT_ALLOWED
    ngx.say(string.format('{"status": %d, "content": "Method Not Allowed: %s"}', ngx.HTTP_NOT_ALLOWED, method))
end
