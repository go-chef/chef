# Legacy resource definition
define :write_a_file, path: nil do
  path = params[:path] || name

  file path
end
