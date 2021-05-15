use std::path::PathBuf;

use crate::config::TargetConfig;

use super::node::{Children, Node, NodeConfig, NodeDefaults};

/// Structure that represents a Pilgo repository of dotfiles.
#[derive(Debug, PartialEq)]
pub struct Tree {
	root: Node,
}

impl Tree {
	pub fn new() -> Tree {
		Tree {
			root: Node::Branch(Children::new()),
		}
	}

	pub fn with(
		mut self,
		path: PathBuf,
		target_config: TargetConfig,
		defaults: NodeDefaults,
	) -> Self {
		self.root.insert(
			path,
			NodeConfig {
				target_config,
				defaults,
				parent_path: PathBuf::new(), // empty parent path for root
			},
		);

		self
	}
}

#[cfg(test)]
mod tests {
	use super::*;

	use std::path::Path;

	use maplit::btreemap;
	use pretty_assertions::assert_eq;

	use crate::config::base_dir::BaseDir;

	#[test]
	fn inserts_single_item() {
		let want = Tree {
			root: Node::Branch(btreemap! {
				Path::new("foo").into() => Node::Leaf {
					target: Path::new("src/foo").into(),
					link: (BaseDir::Config, Path::new("foo").into()),
				},
			}),
		};

		let got = Tree::new().with(
			Path::new("foo").into(),
			TargetConfig::default(),
			(Path::new("src"), BaseDir::Config),
		);

		assert_eq!(got, want);
	}

	#[test]
	fn inserts_nested_item() {
		let want = Tree {
			root: Node::Branch(btreemap! {
				Path::new("foo").into() => Node::Branch(btreemap!{
					Path::new("bar").into() => Node::Leaf {
						target: Path::new("src/foo/bar").into(),
						link: (BaseDir::Config, Path::new("foo/bar").into()),
					},
				}),
			}),
		};

		let got = Tree::new().with(
			Path::new("foo/bar").into(),
			TargetConfig::default(),
			(Path::new("src"), BaseDir::Config),
		);

		assert_eq!(got, want);
	}

	#[test]
	fn override_link() {
		let want = Tree {
			root: Node::Branch(btreemap! {
				Path::new("foo").into() => Node::Leaf {
					target: Path::new("src/foo").into(),
					link: (BaseDir::Config, Path::new("custom").into()),
				},
			}),
		};

		let got = Tree::new().with(
			Path::new("foo").into(),
			TargetConfig {
				link: Some(Path::new("custom").into()),
				..TargetConfig::default()
			},
			(Path::new("src"), BaseDir::Config),
		);

		assert_eq!(got, want);
	}

	#[test]
	fn override_base_dir() {
		let want = Tree {
			root: Node::Branch(btreemap! {
				Path::new("foo").into() => Node::Leaf {
					target: Path::new("src/foo").into(),
					link: (BaseDir::Home, Path::new("foo").into()),
				},
			}),
		};

		let got = Tree::new().with(
			Path::new("foo").into(),
			TargetConfig {
				base_dir: Some(BaseDir::Home),
				..TargetConfig::default()
			},
			(Path::new("src"), BaseDir::Config),
		);

		assert_eq!(got, want);
	}
}
