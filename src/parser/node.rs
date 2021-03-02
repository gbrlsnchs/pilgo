use std::{
	collections::BTreeMap,
	ffi::OsStr,
	path::{Path, PathBuf},
};

use crate::config::TargetConfig;

pub type Children = BTreeMap<PathBuf, Node>;

/// This represents a node in Pilgo's tree. A node can be either simply a branch node, that is, a
/// node that contains one or more children, or a leaf node, finally containing all metadata needed
/// to create symlinks.
#[derive(Debug, PartialEq)]
pub enum Node {
	Branch(Children),
	Leaf { target: PathBuf, link: PathBuf },
}

impl Node {
	pub fn insert<'a>(&mut self, path: PathBuf, config: NodeConfig<'a>) {
		if let Self::Branch(children) = self {
			let segments: Vec<&OsStr> = path.iter().collect();

			if let Some((key, rest)) = segments.split_first() {
				let key = *key;

				if rest.is_empty() {
					let (src_dir, dest_dir) = config.defaults;

					children.insert(
						key.into(),
						Node::Leaf {
							target: Path::new(src_dir).join(key),
							link: Path::new(dest_dir).join(key),
						},
					);
				} else {
					let mut new_node = Node::Branch(Children::new());
					let path = rest.iter().collect();

					new_node.insert(path, config);
					children.insert(key.into(), new_node);
				}
			}
		} else {
			// TODO: This method should probably return an error when trying to insert in a leaf
			// node.
			unimplemented!();
		}
	}
}

pub struct NodeConfig<'a> {
	pub target_config: TargetConfig,
	pub defaults: (&'a Path, &'a Path),
}
