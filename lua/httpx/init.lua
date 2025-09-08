local M = {}

M.config = {
  binary_path = "httpx",
  timeout = 30,
  response_window = {
    width = 0.8,
    height = 0.8,
    border = "rounded",
  },
  keymaps = {
    run_request = "<leader>hr",
    run_all = "<leader>ha",
    show_variables = "<leader>hv",
    select_request = "<leader>hs",
  },
}

function M.setup(opts)
  M.config = vim.tbl_deep_extend("force", M.config, opts or {})
  
  local parser = require("httpx.parser")
  local executor = require("httpx.executor")
  local ui = require("httpx.ui")
  
  executor.setup(M.config)
  ui.setup(M.config)
  
  local function run_request_under_cursor()
    local request = parser.get_request_under_cursor()
    if not request then
      vim.notify("No HTTP request found under cursor", vim.log.levels.WARN)
      return
    end
    
    executor.run_request(request, function(response)
      ui.show_response(response)
    end)
  end
  
  local function run_all_requests()
    local requests = parser.get_all_requests()
    if #requests == 0 then
      vim.notify("No HTTP requests found in file", vim.log.levels.WARN)
      return
    end
    
    executor.run_all_requests(requests, function(responses)
      ui.show_responses(responses)
    end)
  end
  
  local function show_variables()
    local variables = parser.get_variables()
    ui.show_variables(variables)
  end
  
  local function select_and_run_request()
    local requests = parser.get_all_requests()
    if #requests == 0 then
      vim.notify("No HTTP requests found in file", vim.log.levels.WARN)
      return
    end
    
    ui.select_request(requests, function(request)
      if request then
        executor.run_request(request, function(response)
          ui.show_response(response)
        end)
      end
    end)
  end
  
  vim.api.nvim_create_user_command("HttpxRunRequest", run_request_under_cursor, {})
  vim.api.nvim_create_user_command("HttpxRunAll", run_all_requests, {})
  vim.api.nvim_create_user_command("HttpxShowVariables", show_variables, {})
  vim.api.nvim_create_user_command("HttpxSelectRequest", select_and_run_request, {})
  
  vim.api.nvim_create_autocmd("FileType", {
    pattern = "http",
    callback = function(args)
      local buf = args.buf
      vim.keymap.set("n", M.config.keymaps.run_request, run_request_under_cursor, { buffer = buf, desc = "Run HTTP request under cursor" })
      vim.keymap.set("n", M.config.keymaps.run_all, run_all_requests, { buffer = buf, desc = "Run all HTTP requests" })
      vim.keymap.set("n", M.config.keymaps.show_variables, show_variables, { buffer = buf, desc = "Show HTTP variables" })
      vim.keymap.set("n", M.config.keymaps.select_request, select_and_run_request, { buffer = buf, desc = "Select and run HTTP request" })
    end,
  })
end

return M