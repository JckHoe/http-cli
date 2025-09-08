local M = {}

local config = {}
local response_buf = nil
local response_win = nil

function M.setup(opts)
  config = opts
end

local function create_floating_window(title)
  local width = math.floor(vim.o.columns * (config.response_window.width or 0.8))
  local height = math.floor(vim.o.lines * (config.response_window.height or 0.8))
  local row = math.floor((vim.o.lines - height) / 2)
  local col = math.floor((vim.o.columns - width) / 2)
  
  local buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_buf_set_option(buf, "buftype", "nofile")
  vim.api.nvim_buf_set_option(buf, "bufhidden", "wipe")
  vim.api.nvim_buf_set_option(buf, "filetype", "httpx-response")
  
  local win_opts = {
    relative = "editor",
    width = width,
    height = height,
    row = row,
    col = col,
    style = "minimal",
    border = config.response_window.border or "rounded",
    title = title or " HTTP Response ",
    title_pos = "center",
  }
  
  local win = vim.api.nvim_open_win(buf, true, win_opts)
  
  vim.api.nvim_buf_set_keymap(buf, "n", "q", ":close<CR>", { noremap = true, silent = true })
  vim.api.nvim_buf_set_keymap(buf, "n", "<Esc>", ":close<CR>", { noremap = true, silent = true })
  
  vim.api.nvim_win_set_option(win, "number", false)
  vim.api.nvim_win_set_option(win, "relativenumber", false)
  vim.api.nvim_win_set_option(win, "cursorline", true)
  vim.api.nvim_win_set_option(win, "wrap", true)
  
  return buf, win
end

local function format_response_lines(response)
  local lines = {}
  
  if response.error then
    table.insert(lines, "═══ ERROR ═══")
    table.insert(lines, "")
    table.insert(lines, response.error)
  else
    if response.request then
      table.insert(lines, "═══ REQUEST ═══")
      table.insert(lines, "")
      table.insert(lines, response.request.method .. " " .. response.request.url)
      if response.request.name and response.request.name ~= "" then
        table.insert(lines, "Name: " .. response.request.name)
      end
      table.insert(lines, "")
    end
    
    table.insert(lines, "═══ RESPONSE ═══")
    table.insert(lines, "")
    
    if response.status ~= "" then
      table.insert(lines, "Status: " .. response.status)
    end
    
    if response.duration ~= "" then
      table.insert(lines, "Duration: " .. response.duration)
    end
    
    if #response.headers > 0 then
      table.insert(lines, "")
      table.insert(lines, "Headers:")
      for _, header in ipairs(response.headers) do
        table.insert(lines, "  " .. header.key .. ": " .. header.value)
      end
    end
    
    if response.body and response.body ~= "" then
      table.insert(lines, "")
      table.insert(lines, "Body:")
      table.insert(lines, "")
      
      local body_lines = vim.split(response.body, "\n")
      for _, line in ipairs(body_lines) do
        table.insert(lines, line)
      end
    end
  end
  
  return lines
end

local function apply_syntax_highlighting(buf)
  vim.api.nvim_buf_call(buf, function()
    vim.cmd([[
      syntax match CassieHttpSeparator /^═\+.*═\+$/
      syntax match CassieHttpHeader /^[A-Za-z-]\+:/ nextgroup=CassieHttpHeaderValue
      syntax match CassieHttpHeaderValue /.*/ contained
      syntax match CassieHttpStatus /^Status:.*$/
      syntax match CassieHttpDuration /^Duration:.*$/
      syntax match CassieHttpMethod /^\(GET\|POST\|PUT\|DELETE\|PATCH\|HEAD\|OPTIONS\|TRACE\|CONNECT\)/
      syntax match CassieHttpUrl /https\?:\/\/[^ ]*/
      syntax match CassieHttpError /^Error:.*$/
      syntax region CassieHttpJson start=/{/ end=/}/ contains=CassieHttpJsonKey,CassieHttpJsonString,CassieHttpJsonNumber,CassieHttpJsonBoolean,CassieHttpJsonNull fold
      syntax region CassieHttpJson start=/\[/ end=/\]/ contains=CassieHttpJsonKey,CassieHttpJsonString,CassieHttpJsonNumber,CassieHttpJsonBoolean,CassieHttpJsonNull fold
      syntax match CassieHttpJsonKey /"[^"]*":\ze/
      syntax region CassieHttpJsonString start=/"/ skip=/\\"/ end=/"/
      syntax match CassieHttpJsonNumber /\<-\?\d\+\(\.\d\+\)\?\([eE][+-]\?\d\+\)\?\>/
      syntax keyword CassieHttpJsonBoolean true false
      syntax keyword CassieHttpJsonNull null
      
      highlight link CassieHttpSeparator Title
      highlight link CassieHttpHeader Label
      highlight link CassieHttpHeaderValue String
      highlight link CassieHttpStatus Function
      highlight link CassieHttpDuration Number
      highlight link CassieHttpMethod Keyword
      highlight link CassieHttpUrl Underlined
      highlight link CassieHttpError ErrorMsg
      highlight link CassieHttpJsonKey Identifier
      highlight link CassieHttpJsonString String
      highlight link CassieHttpJsonNumber Number
      highlight link CassieHttpJsonBoolean Boolean
      highlight link CassieHttpJsonNull Constant
    ]])
  end)
end

function M.show_response(response)
  response_buf, response_win = create_floating_window(" HTTP Response ")
  
  local lines = format_response_lines(response)
  vim.api.nvim_buf_set_lines(response_buf, 0, -1, false, lines)
  vim.api.nvim_buf_set_option(response_buf, "modifiable", false)
  
  apply_syntax_highlighting(response_buf)
end

function M.show_responses(responses)
  response_buf, response_win = create_floating_window(" HTTP Test Results ")
  
  local lines = {}
  
  if responses.error and responses.error ~= "" then
    table.insert(lines, "═══ ERROR ═══")
    table.insert(lines, "")
    table.insert(lines, responses.error)
  else
    vim.list_extend(lines, vim.split(responses.output, "\n"))
  end
  
  vim.api.nvim_buf_set_lines(response_buf, 0, -1, false, lines)
  vim.api.nvim_buf_set_option(response_buf, "modifiable", false)
  
  apply_syntax_highlighting(response_buf)
end

function M.show_variables(variables)
  local buf, win = create_floating_window(" HTTP Variables ")
  
  local lines = {}
  table.insert(lines, "═══ VARIABLES ═══")
  table.insert(lines, "")
  
  if next(variables) then
    for name, value in pairs(variables) do
      table.insert(lines, string.format("@%s = %s", name, value))
    end
  else
    table.insert(lines, "No variables defined")
  end
  
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)
  vim.api.nvim_buf_set_option(buf, "modifiable", false)
  
  apply_syntax_highlighting(buf)
end

function M.select_request(requests, callback)
  local items = {}
  
  for i, request in ipairs(requests) do
    local label = string.format("%d. %s %s", i, request.method, request.url)
    if request.name and request.name ~= "" then
      label = string.format("%d. [%s] %s %s", i, request.name, request.method, request.url)
    end
    table.insert(items, label)
  end
  
  vim.ui.select(items, {
    prompt = "Select HTTP request to run:",
    format_item = function(item)
      return item
    end,
  }, function(choice, idx)
    if choice and idx then
      callback(requests[idx])
    else
      callback(nil)
    end
  end)
end

return M