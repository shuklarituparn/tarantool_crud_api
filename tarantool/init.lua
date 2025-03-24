box.cfg {
    listen = 3301
}

box.once('init_kv_store', function()
    local kv_store = box.schema.space.create('kv_store', {if_not_exists = true})
    
    kv_store:format({
        {name = 'key', type = 'string'},
        {name = 'value', type = 'any'}
    })
    
    kv_store:create_index('primary', {
        type = 'hash',
        parts = {{field = 1, type = 'string'}},
        if_not_exists = true
    })
    
    if not box.schema.user.exists('kv_user') then
        box.schema.user.create('kv_user', {
            password = 'Hellopassword123'
        })
        
        box.schema.user.grant('kv_user', 'read,write,execute', 'universe')
    end
end)
