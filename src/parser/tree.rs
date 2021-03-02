use std::path::{Path, PathBuf};

use super::node::{Node, NodeConfig};
use crate::config::TargetConfig;

/// Structure that represents a Pilgo repository of dotfiles.
#[derive(Debug, PartialEq)]
pub struct Tree<'a> {
	base_dirs: (&'a Path, &'a Path), // src and dest
	root: Node,
}

impl<'a> Tree<'a> {
	pub fn insert(&mut self, path: PathBuf, target_config: TargetConfig) {
		self.root.insert(
			path,
			NodeConfig {
				target_config,
				defaults: self.base_dirs,
			},
		);
	}
}

#[cfg(test)]
mod tests {
	use super::*;

	use maplit::btreemap;
	use pretty_assertions::assert_eq;

	use crate::parser::node::Children;

	#[test]
	fn inserts_single_item() {
		let want = Tree {
			base_dirs: (Path::new("src"), Path::new("dest")),
			root: Node::Branch(btreemap! {
				Path::new("foo").into() => Node::Leaf {
					target: Path::new("src/foo").into(),
					link: Path::new("dest/foo").into(),
				},
			}),
		};

		let mut got = Tree {
			base_dirs: (Path::new("src"), Path::new("dest")),
			root: Node::Branch(Children::new()),
		};
		got.insert(Path::new("foo").into(), TargetConfig::default());

		assert_eq!(got, want);
	}
}
