local M = {}

local function get_lines()
  return vim.api.nvim_buf_get_lines(0, 0, -1, false)
end

local function is_separator(line)
  return line:match("^###%s*(.*)$")
end

local function is_request_line(line)
  return line:match("^(%u+)%s+(.+)") or line:match("^https?://")
end

local function is_header_line(line)
  return line:match("^([%w%-]+):%s*(.+)")
end

local function is_variable_line(line)
  return line:match("^@([%w_]+)%s*=%s*(.+)")
end

local function is_comment(line)
  return line:match("^#[^#]") or line:match("^//")
end

function M.get_request_under_cursor()
  local cursor_line = vim.fn.line(".")
  local lines = get_lines()
  
  local request_start = cursor_line
  local request_end = cursor_line
  
  for i = cursor_line, 1, -1 do
    if is_separator(lines[i]) then
      request_start = i + 1
      break
    elseif i == 1 then
      request_start = 1
    end
  end
  
  for i = cursor_line, #lines do
    if is_separator(lines[i]) and i > cursor_line then
      request_end = i - 1
      break
    elseif i == #lines then
      request_end = #lines
    end
  end
  
  local request = {
    start_line = request_start,
    end_line = request_end,
    method = nil,
    url = nil,
    headers = {},
    body = "",
    name = "",
  }
  
  local in_body = false
  local body_lines = {}
  
  if request_start > 1 and is_separator(lines[request_start - 1]) then
    local separator_match = lines[request_start - 1]:match("^###%s*(.*)$")
    request.name = separator_match or ""
  end
  
  for i = request_start, request_end do
    local line = lines[i]
    
    if is_comment(line) then
      goto continue
    end
    
    if not in_body then
      if not request.method and is_request_line(line) then
        local method, url = line:match("^(%u+)%s+(.+)")
        if method and url then
          request.method = method
          request.url = url
        elseif line:match("^https?://") then
          request.method = "GET"
          request.url = line:match("^%s*(.-)%s*$")
        end
      elseif request.method and is_header_line(line) then
        local key, value = line:match("^([%w%-]+):%s*(.+)")
        if key and value then
          table.insert(request.headers, { key = key, value = value })
        end
      elseif request.method and line:match("^%s*$") then
        in_body = true
      end
    else
      table.insert(body_lines, line)
    end
    
    ::continue::
  end
  
  if #body_lines > 0 then
    while #body_lines > 0 and body_lines[#body_lines]:match("^%s*$") do
      table.remove(body_lines)
    end
    request.body = table.concat(body_lines, "\n")
  end
  
  if not request.method then
    return nil
  end
  
  return request
end

function M.get_all_requests()
  local lines = get_lines()
  local requests = {}
  local current_request = nil
  local in_body = false
  local body_lines = {}
  
  for i, line in ipairs(lines) do
    if is_separator(line) then
      if current_request and current_request.method then
        if #body_lines > 0 then
          while #body_lines > 0 and body_lines[#body_lines]:match("^%s*$") do
            table.remove(body_lines)
          end
          current_request.body = table.concat(body_lines, "\n")
        end
        table.insert(requests, current_request)
      end
      
      local separator_match = line:match("^###%s*(.*)$")
      current_request = {
        start_line = i + 1,
        method = nil,
        url = nil,
        headers = {},
        body = "",
        name = separator_match or "",
      }
      in_body = false
      body_lines = {}
    elseif is_comment(line) then
    elseif not current_request then
      if not line:match("^%s*$") then
        current_request = {
          start_line = i,
          method = nil,
          url = nil,
          headers = {},
          body = "",
          name = "",
        }
      end
    end
    
    if current_request then
      if not in_body then
        if not current_request.method and is_request_line(line) then
          local method, url = line:match("^(%u+)%s+(.+)")
          if method and url then
            current_request.method = method
            current_request.url = url
          elseif line:match("^https?://") then
            current_request.method = "GET"
            current_request.url = line:match("^%s*(.-)%s*$")
          end
        elseif current_request.method and is_header_line(line) then
          local key, value = line:match("^([%w%-]+):%s*(.+)")
          if key and value then
            table.insert(current_request.headers, { key = key, value = value })
          end
        elseif current_request.method and line:match("^%s*$") then
          in_body = true
        end
      else
        if not is_separator(line) then
          table.insert(body_lines, line)
        end
      end
    end
  end
  
  if current_request and current_request.method then
    if #body_lines > 0 then
      while #body_lines > 0 and body_lines[#body_lines]:match("^%s*$") do
        table.remove(body_lines)
      end
      current_request.body = table.concat(body_lines, "\n")
    end
    current_request.end_line = #lines
    table.insert(requests, current_request)
  end
  
  return requests
end

function M.get_variables()
  local lines = get_lines()
  local variables = {}
  
  for _, line in ipairs(lines) do
    if is_variable_line(line) then
      local name, value = line:match("^@([%w_]+)%s*=%s*(.+)")
      if name and value then
        variables[name] = value
      end
    end
  end
  
  return variables
end

return M