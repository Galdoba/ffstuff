package copyprocess

func (cas *copyActionState) validatePaths() error {
	//destinations checks
	// destInfo, errD := os.Stat(cas.dest)
	// if errD != nil {
	// 	return errors.New("Destination: " + errD.Error())
	// }
	// if !destInfo.IsDir() {
	// 	return errors.New("Destination is not a directory: " + destInfo.Name())
	// }
	// for _, source := range cas.sources {
	// 	srcInfo, errS := os.Stat(source)
	// 	if errS != nil {
	// 		return fmt.Errorf("source: stats gathering failed: %v", errS)
	// 	}
	// 	if !srcInfo.Mode().IsRegular() { // cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
	// 		return fmt.Errorf("source: non-regular file: " + srcInfo.Name() + " (" + srcInfo.Mode().String() + ")")
	// 	}
	// }
	return nil
}
