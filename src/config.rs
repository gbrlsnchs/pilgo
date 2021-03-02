use std::{collections::BTreeMap, path::PathBuf};

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
	/// This is an optional property that renames the link name for the target.
	pub link: Option<PathBuf>,
}
