use std::{collections::BTreeMap, path::PathBuf};

use self::base_dir::BaseDir;

pub mod base_dir;

/// The configuration is the source of truth for Pilgo. It has information about targets and their
/// individual settings.
pub struct Config {
	/// All targets on which Pilgo should operate.
	pub targets: Vec<PathBuf>,
	/// Individual options for each target, if any.
	pub options: BTreeMap<PathBuf, TargetConfig>,
}

/// Optional, individual target configuration.
#[derive(Default)]
pub struct TargetConfig {
	/// Custom link name.
	pub link: Option<PathBuf>,
	/// Custom base directory for the link.
	pub base_dir: Option<BaseDir>,
}
