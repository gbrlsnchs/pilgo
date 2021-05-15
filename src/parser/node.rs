use std::{
	collections::BTreeMap,
	ffi::OsStr,
	path::{Path, PathBuf},
};

use crate::config::{base_dir::BaseDir, TargetConfig};

pub type Children = BTreeMap<PathBuf, Node>;

/// This represents a node in Pilgo's tree. A node can be either simply a branch node, that is, a
/// node that contains one or more children, or a leaf node, finally containing all metadata needed
/// to create symlinks.
#[derive(Debug, PartialEq)]
pub enum Node {
	Branch(Children),
	Leaf {
		target: PathBuf,
		link: (BaseDir, PathBuf),
	},
}

impl Node {
	fn new_branch(path: PathBuf, config: NodeConfig) -> Self {
		let mut node = Node::Branch(Children::new());

		node.insert(path, config);

		node
	}

	fn new_leaf<T, L>(
		src_dir: &Path,
		dest_dir: BaseDir,
		parent_path: PathBuf,
		target: T,
		link: L,
	) -> Self
	where
		T: AsRef<Path>,
		L: AsRef<Path>,
	{
		Node::Leaf {
			target: Path::new(src_dir).join(&parent_path).join(target),
			link: (dest_dir, Path::new(&parent_path).join(link)),
		}
	}

	pub fn insert(&mut self, path: PathBuf, mut config: NodeConfig) {
		if let Self::Branch(children) = self {
			let segments: Vec<&OsStr> = path.iter().collect();

			if let Some((key, rest)) = segments.split_first() {
				let key = *key;

				if rest.is_empty() {
					let (src_dir, dest_dir) = config.defaults;
					let NodeConfig {
						target_config,
						parent_path,
						..
					} = config;

					let dest_dir = target_config.base_dir.unwrap_or(dest_dir);

					let link = target_config.link.unwrap_or_else(|| key.into());

					let leaf = Node::new_leaf(src_dir, dest_dir, parent_path, key, link);
					children.insert(key.into(), leaf);
				} else {
					let path = rest.iter().collect();

					config.parent_path = Path::new(key).join(config.parent_path);

					let branch = Node::new_branch(path, config);
					children.insert(key.into(), branch);
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
	pub defaults: NodeDefaults<'a>,
	pub parent_path: PathBuf,
}

pub type NodeDefaults<'a> = (&'a Path, BaseDir);
