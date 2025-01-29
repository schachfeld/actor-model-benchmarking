defmodule Coinbaseanalyser.MixProject do
  use Mix.Project

  def project do
    [
      app: :coinbaseanalyser,
      version: "0.1.0",
      elixir: "~> 1.18-dev",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger],
      mod: {FileReader, []}
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      {:poison, "~> 6.0"}
      # {:dep_from_git, git: "https://github.com/elixir-lang/my_dep.git", tag: "0.1.0"}
    ]
  end
end
