
-- RoleId以1结尾的玩家拒绝登陆
function RoleHook(roleid)
    if roleid:sub(-1, -1) == '1' then
        return false
    end
    return true
end