use std::path::{Path, PathBuf};

use super::node::{Children, Node, NodeConfig};
use crate::config::TargetConfig;

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
		defaults: (&Path, &Path),
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

	use maplit::btreemap;
	use pretty_assertions::assert_eq;

	#[test]
	fn inserts_single_item() {
		let want = Tree {
			root: Node::Branch(btreemap! {
				Path::new("foo").into() => Node::Leaf {
					target: Path::new("src/foo").into(),
					link: Path::new("dest/foo").into(),
				},
			}),
		};

		let got = Tree::new().with(
			Path::new("foo").into(),
			TargetConfig::default(),
			(Path::new("src"), Path::new("dest")),
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
						link: Path::new("dest/foo/bar").into(),
					},
				}),
			}),
		};

		let got = Tree::new().with(
			Path::new("foo/bar").into(),
			TargetConfig::default(),
			(Path::new("src"), Path::new("dest")),
		);

		assert_eq!(got, want);
	}
}
