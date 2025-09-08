local M = {}

local config = {}

function M.setup(opts)
  config = opts
end

local function create_temp_file(content)
  local tmpfile = vim.fn.tempname() .. ".http"
  local file = io.open(tmpfile, "w")
  if not file then
    return nil
  end
  file:write(content)
  file:close()
  return tmpfile
end

local function format_request(request)
  local lines = {}
  
  if request.name and request.name ~= "" then
    table.insert(lines, "### " .. request.name)
  else
    table.insert(lines, "###")
  end
  
  table.insert(lines, request.method .. " " .. request.url)
  
  for _, header in ipairs(request.headers) do
    table.insert(lines, header.key .. ": " .. header.value)
  end
  
  if request.body and request.body ~= "" then
    table.insert(lines, "")
    table.insert(lines, request.body)
  end
  
  return table.concat(lines, "\n")
end

local function format_all_requests(requests)
  local parser = require("httpx.parser")
  local variables = parser.get_variables()
  local lines = {}
  
  for name, value in pairs(variables) do
    table.insert(lines, "@" .. name .. " = " .. value)
  end
  
  if next(variables) then
    table.insert(lines, "")
  end
  
  for i, request in ipairs(requests) do
    table.insert(lines, format_request(request))
    if i < #requests then
      table.insert(lines, "")
    end
  end
  
  return table.concat(lines, "\n")
end

local function parse_response(output, stderr)
  local response = {
    status = "",
    duration = "",
    headers = {},
    body = "",
    error = nil,
  }
  
  if stderr and stderr ~= "" then
    response.error = stderr
    return response
  end
  
  local lines = vim.split(output, "\n")
  local in_body = false
  local body_lines = {}
  local in_headers = false
  
  for _, line in ipairs(lines) do
    if line:match("^Status:") then
      response.status = line:match("^Status:%s*(.+)")
    elseif line:match("^Duration:") then
      response.duration = line:match("^Duration:%s*(.+)")
    elseif line:match("^Error:") then
      response.error = line:match("^Error:%s*(.+)")
    elseif line:match("^Headers:") then
      in_headers = true
      in_body = false
    elseif line:match("^Body:") then
      in_headers = false
      in_body = true
    elseif in_headers and line:match("^%s+") then
      local key, value = line:match("^%s+([^:]+):%s*(.+)")
      if key and value then
        table.insert(response.headers, { key = key, value = value })
      end
    elseif in_body then
      table.insert(body_lines, line)
    end
  end
  
  if #body_lines > 0 then
    response.body = table.concat(body_lines, "\n")
  end
  
  return response
end

function M.run_request(request, callback)
  local content = format_request(request)
  local tmpfile = create_temp_file(content)
  
  if not tmpfile then
    vim.notify("Failed to create temporary file", vim.log.levels.ERROR)
    return
  end
  
  local cmd = {
    config.binary_path or "httpx",
    "run",
    tmpfile,
    "--request", "1",
  }
  
  if config.timeout then
    table.insert(cmd, "--timeout")
    table.insert(cmd, tostring(config.timeout) .. "s")
  end
  
  local output = {}
  local stderr = {}
  
  vim.fn.jobstart(cmd, {
    stdout_buffered = true,
    stderr_buffered = true,
    on_stdout = function(_, data)
      if data then
        vim.list_extend(output, data)
      end
    end,
    on_stderr = function(_, data)
      if data then
        vim.list_extend(stderr, data)
      end
    end,
    on_exit = function(_, exit_code)
      vim.schedule(function()
        vim.fn.delete(tmpfile)
        
        local response = parse_response(
          table.concat(output, "\n"),
          table.concat(stderr, "\n")
        )
        
        response.request = request
        response.exit_code = exit_code
        
        if callback then
          callback(response)
        end
      end)
    end,
  })
end

function M.run_all_requests(requests, callback)
  local content = format_all_requests(requests)
  local tmpfile = create_temp_file(content)
  
  if not tmpfile then
    vim.notify("Failed to create temporary file", vim.log.levels.ERROR)
    return
  end
  
  local cmd = {
    config.binary_path or "httpx",
    "test",
    tmpfile,
  }
  
  if config.timeout then
    table.insert(cmd, "--timeout")
    table.insert(cmd, tostring(config.timeout) .. "s")
  end
  
  local output = {}
  local stderr = {}
  
  vim.fn.jobstart(cmd, {
    stdout_buffered = true,
    stderr_buffered = true,
    on_stdout = function(_, data)
      if data then
        vim.list_extend(output, data)
      end
    end,
    on_stderr = function(_, data)
      if data then
        vim.list_extend(stderr, data)
      end
    end,
    on_exit = function(_, exit_code)
      vim.schedule(function()
        vim.fn.delete(tmpfile)
        
        local responses = {
          output = table.concat(output, "\n"),
          error = table.concat(stderr, "\n"),
          exit_code = exit_code,
          requests = requests,
        }
        
        if callback then
          callback(responses)
        end
      end)
    end,
  })
end

return M