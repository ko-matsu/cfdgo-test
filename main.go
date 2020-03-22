package main

import (
    "fmt"
    "time"
    "github.com/pkg/profile"
)
import cfd "github.com/ko-matsu/cfd-go"
//     "encoding/hex"

func main() {
    defer profile.Start(profile.ProfilePath(".")).Stop()

    fmt.Println("start")
    TestCfdCreateLargeTransaction()
    // CreateBasicTransaction()
    // TestCfdCreateLargeBlindTransaction()
    fmt.Println("end. finishing...")
    time.Sleep(10 * time.Second) // GC wait
    fmt.Println("finish")
}

const maxnum = 102400

func TestCfdCreateLargeTransaction() {
	handle := uintptr(0)
	sequence := (uint32)(cfd.KCfdSequenceLockTimeDisable)
	networkType := (int)(cfd.KCfdNetworkLiquidv1)
	maxTxin := maxnum
	maxTxout := maxnum
	// mnemonic: accuse traffic neglect mechanic sand page cycle tattoo bonus sheriff field top vote outdoor drop
	// master xpriv: tprv8ZgxMBicQKsPe4wfWRxqwNiLTkpLoSvHsKkSP9XXfbV1KVRazP6HNhfDb85caijRwMN8VKFwQu6WeCFjUQGeKk4MTDpHfAahaNHuhUfXhYt
	masterXpriv := "tprv8ZgxMBicQKsPe4wfWRxqwNiLTkpLoSvHsKkSP9XXfbV1KVRazP6HNhfDb85caijRwMN8VKFwQu6WeCFjUQGeKk4MTDpHfAahaNHuhUfXhYt"
	asset := "ef47c42d34de1b06a02212e8061323f50d5f02ceed202f1cb375932aa299f751"

	createTxHandle, err := cfd.CfdGoInitializeConfidentialTransaction(uint32(2), uint32(0))
	if err != nil {
		return
	}
	defer cfd.CfdGoFreeTransactionHandle(createTxHandle)

	for i := 1; i <= maxTxin; i++ {
		txid := "00000000000000000000000000000000000000000000000000000000" + fmt.Sprintf("%08x", i)
		err = cfd.CfdGoAddTxInput(
			createTxHandle, txid, 0, sequence)
		if err != nil {
			fmt.Printf("CfdGoAddConfidentialTxIn fail[%s] idx[%d]\n", err.Error(), i)
			return
		}

		if i % 128 == 0 {
			fmt.Print(" - txin: " + txid + "\n")
		}
	}

	descriptor := "pkh(" + masterXpriv + "/0/*)"
	for i := 1; i < maxTxout; i++ {
		bip32DerivationPath := fmt.Sprintf("%d", i)
		descriptorDataList, _, err := cfd.CfdGoParseDescriptor(handle , descriptor, networkType, bip32DerivationPath)
		if err != nil {
			fmt.Printf("CfdGoParseDescriptor fail[%s] idx[%d]\n", err.Error(), i)
			return
		}
	
		err = cfd.CfdGoAddConfidentialTxOutput(
			createTxHandle, asset, int64(10000000),
			descriptorDataList[0].Address)
		if err != nil {
			fmt.Print("CfdGoAddConfidentialTxOut fail[%s] idx[%d]\n", err.Error(), i)
			return
		}

		if i % 128 == 0 {
			fmt.Printf(" - txout: %d\n", i)
		}
	}

	bip32DerivationPath := fmt.Sprintf("%d", maxTxout)
	descriptorDataList, _, err := cfd.CfdGoParseDescriptor(handle , descriptor, networkType, bip32DerivationPath)
	if err != nil {
		fmt.Print("CfdGoParseDescriptor fail[%s] idx[%d]\n", err.Error(), maxTxout)
		return
	}
	err = cfd.CfdGoAddConfidentialTxOutput(
		createTxHandle, asset, int64(9900000),
		descriptorDataList[0].Address)
	if err != nil {
		fmt.Printf("CfdGoAddConfidentialTxOut fail[%s] idx[%d]\n", err.Error(), maxTxout)
		return
	}

	err = cfd.CfdGoAddConfidentialTxOutputFee(
		createTxHandle, asset, int64(100000))
	if err != nil {
		fmt.Printf("CfdGoAddConfidentialTxOut fail[%s] idx[fee]\n", err.Error())
		return
	}

	txHex, err := cfd.CfdGoFinalizeTransaction(createTxHandle)
	if err != nil {
		fmt.Printf("CfdGoFinalizeTransaction fail[%s]\n", err.Error())
		return
	}

	if err != nil {
		errMsg, _ := cfd.CfdGoGetLastErrorMessage(handle)
		fmt.Print("[error message] " + errMsg + "\n")
	}

	fmt.Printf("txHex = %s \n", txHex)
	fmt.Print("TestCfdCreateLargeTransaction test done.\n")
}

func CreateBasicTransaction() (txHex string, err error) {
	handle := uintptr(0)
	sequence := (uint32)(cfd.KCfdSequenceLockTimeDisable)
	networkType := (int)(cfd.KCfdNetworkLiquidv1)
	maxTxin := maxnum
	maxTxout := maxnum
	// mnemonic: accuse traffic neglect mechanic sand page cycle tattoo bonus sheriff field top vote outdoor drop
	// master xpriv: tprv8ZgxMBicQKsPe4wfWRxqwNiLTkpLoSvHsKkSP9XXfbV1KVRazP6HNhfDb85caijRwMN8VKFwQu6WeCFjUQGeKk4MTDpHfAahaNHuhUfXhYt
	masterXpriv := "tprv8ZgxMBicQKsPe4wfWRxqwNiLTkpLoSvHsKkSP9XXfbV1KVRazP6HNhfDb85caijRwMN8VKFwQu6WeCFjUQGeKk4MTDpHfAahaNHuhUfXhYt"
	asset := "ef47c42d34de1b06a02212e8061323f50d5f02ceed202f1cb375932aa299f751"

	txHex, err = cfd.CfdGoInitializeConfidentialTx(handle, uint32(2), uint32(0))
	if err != nil {
		return
	}

	for i := 1; i <= maxTxin; i++ {
		txid := "00000000000000000000000000000000000000000000000000000000" + fmt.Sprintf("%08x", i)
		txHex, err = cfd.CfdGoAddConfidentialTxIn(
			handle, txHex, txid, 0, sequence)
		if err != nil {
			fmt.Print("CfdGoAddConfidentialTxIn fail[%s] idx[%d]", err.Error(), i)
			return
		}

		if i % 128 == 0 {
			fmt.Print(" - txin: " + txid + "\n")
		}
	}

	descriptor := "pkh(" + masterXpriv + "/0/*)"
	for i := 1; i < maxTxout; i++ {
		bip32DerivationPath := fmt.Sprintf("%d", i)
		descriptorDataList, _, err := cfd.CfdGoParseDescriptor(handle , descriptor, networkType, bip32DerivationPath)
		if err != nil {
			fmt.Printf("CfdGoParseDescriptor fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}
	
		txHex, err = cfd.CfdGoAddConfidentialTxOut(
			handle, txHex, asset, int64(10000000), "",
			descriptorDataList[0].Address, "", "")
		if err != nil {
			fmt.Printf("CfdGoAddConfidentialTxOut fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}

		if i % 128 == 0 {
			fmt.Printf(" - txout: %d\n", i)
		}
	}
	bip32DerivationPath := fmt.Sprintf("%d", maxTxout)
	descriptorDataList, _, err := cfd.CfdGoParseDescriptor(handle , descriptor, networkType, bip32DerivationPath)
	if err != nil {
		fmt.Print("CfdGoParseDescriptor fail[%s] idx[%d]\n", err.Error(), maxTxout)
		return
	}
	txHex, err = cfd.CfdGoAddConfidentialTxOut(
		handle, txHex, asset, int64(9900000), "",
		descriptorDataList[0].Address, "", "")
	if err != nil {
		fmt.Printf("CfdGoAddConfidentialTxOut fail[%s] idx[%d]\n", err.Error(), maxTxout)
		return
	}

	txHex, err = cfd.CfdGoAddConfidentialTxOut(
		handle, txHex, asset, int64(100000), "", "", "", "")
	if err != nil {
		fmt.Printf("CfdGoAddConfidentialTxOut fail[%s] idx[fee]\n", err.Error())
		return
	}

	return txHex, err
}

func CreateBasicBlindTransaction() (txHex string, err error) {
	handle := uintptr(0)
	sequence := (uint32)(cfd.KCfdSequenceLockTimeDisable)
	networkType := (int)(cfd.KCfdNetworkLiquidv1)
	maxTxin := maxnum
	maxTxout := maxnum
	// mnemonic: accuse traffic neglect mechanic sand page cycle tattoo bonus sheriff field top vote outdoor drop
	// master xpriv: tprv8ZgxMBicQKsPe4wfWRxqwNiLTkpLoSvHsKkSP9XXfbV1KVRazP6HNhfDb85caijRwMN8VKFwQu6WeCFjUQGeKk4MTDpHfAahaNHuhUfXhYt
	masterXpriv := "tprv8ZgxMBicQKsPe4wfWRxqwNiLTkpLoSvHsKkSP9XXfbV1KVRazP6HNhfDb85caijRwMN8VKFwQu6WeCFjUQGeKk4MTDpHfAahaNHuhUfXhYt"
	asset := "ef47c42d34de1b06a02212e8061323f50d5f02ceed202f1cb375932aa299f751"

	txHex, err = cfd.CfdGoInitializeConfidentialTx(handle, uint32(2), uint32(0))
	if err != nil {
		return
	}

	for i := 1; i <= maxTxin; i++ {
		txid := "00000000000000000000000000000000000000000000000000000000" + fmt.Sprintf("%08x", i)
		txHex, err = cfd.CfdGoAddConfidentialTxIn(
			handle, txHex, txid, 0, sequence)
		if err != nil {
			fmt.Print("CfdGoAddConfidentialTxIn fail[%s] idx[%d]", err.Error(), i)
			return
		}
	}

	descriptor := "pkh(" + masterXpriv + "/0/*)"
	for i := 1; i < maxTxout; i++ {
		bip32DerivationPath := fmt.Sprintf("%d", i)
		descriptorDataList, _, err := cfd.CfdGoParseDescriptor(handle , descriptor, networkType, bip32DerivationPath)
		if err != nil {
			fmt.Printf("CfdGoParseDescriptor fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}
	
		txHex, err = cfd.CfdGoAddConfidentialTxOut(
			handle, txHex, asset, int64(10000000), "",
			descriptorDataList[0].Address, "", "")
		if err != nil {
			fmt.Printf("CfdGoAddConfidentialTxOut fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}
	}
	bip32DerivationPath := fmt.Sprintf("%d", maxTxout)
	descriptorDataList, _, err := cfd.CfdGoParseDescriptor(handle , descriptor, networkType, bip32DerivationPath)
	if err != nil {
		fmt.Print("CfdGoParseDescriptor fail[%s] idx[%d]\n", err.Error(), maxTxout)
		return
	}
	txHex, err = cfd.CfdGoAddConfidentialTxOut(
		handle, txHex, asset, int64(9900000), "",
		descriptorDataList[0].Address, "", "")
	if err != nil {
		fmt.Printf("CfdGoAddConfidentialTxOut fail[%s] idx[%d]\n", err.Error(), maxTxout)
		return
	}

	txHex, err = cfd.CfdGoAddConfidentialTxOut(
		handle, txHex, asset, int64(100000), "", "", "", "")
	if err != nil {
		fmt.Printf("CfdGoAddConfidentialTxOut fail[%s] idx[fee]\n", err.Error())
		return
	}


	// blind
	fmt.Print("blinding start.\n")
	blindHandle, err := cfd.CfdGoInitializeBlindTx(handle)
	if err != nil {
		fmt.Print("CfdGoInitializeBlindTx fail[%s]\n", err.Error())
		return "", err
	}
	defer cfd.CfdGoFreeBlindHandle(handle, blindHandle)
	
	emptyBlinder := "0000000000000000000000000000000000000000000000000000000000000000"
	for i := 1; i <= maxTxin; i++ {
		txid := "00000000000000000000000000000000000000000000000000000000" + fmt.Sprintf("%08x", i)
		err = cfd.CfdGoAddBlindTxInData(handle, blindHandle, txid, uint32(0), asset, emptyBlinder, emptyBlinder, int64(10000000), "", "")
		if err != nil {
			fmt.Printf("CfdGoAddBlindTxInData fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}
	}

	var path string
	for i := 1; i <= maxTxout; i++ {
		path = fmt.Sprintf("1/%d", i)
		childKey, err := cfd.CfdGoCreateExtkeyFromParentPath(handle, masterXpriv, path, (int)(cfd.KCfdNetworkTestnet), (int)(cfd.KCfdExtPrivkey))
		if err != nil {
			fmt.Printf("CfdGoCreateExtkeyFromParentPath fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}
		confidentialKey, err := cfd.CfdGoGetPubkeyFromExtkey(handle, childKey, (int)(cfd.KCfdNetworkTestnet))
		if err != nil {
			fmt.Printf("CfdGoGetPubkeyFromExtkey fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}
	
		err = cfd.CfdGoAddBlindTxOutData(handle, blindHandle, uint32(i - 1), confidentialKey)
		if err != nil {
			fmt.Printf("CfdGoAddBlindTxOutData fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}
	}

	txHex, err = cfd.CfdGoFinalizeBlindTx(handle, blindHandle, txHex)
	if err != nil {
		fmt.Printf("CfdGoFinalizeBlindTx fail[%s]\n", err.Error())
		return "", err
	}
	fmt.Print("blinding end.\n")
	return txHex, err
}

func TestCfdBlindTransaction(txHex string, baseTxHex string) (outputTxHex string, err error) {
	handle := uintptr(0)
	maxTxin := maxnum
	maxTxout := maxnum
	masterXpriv := "tprv8ZgxMBicQKsPe4wfWRxqwNiLTkpLoSvHsKkSP9XXfbV1KVRazP6HNhfDb85caijRwMN8VKFwQu6WeCFjUQGeKk4MTDpHfAahaNHuhUfXhYt"

	fmt.Print("CfdGoGetConfidentialTxData start.\n")
	txData, err := cfd.CfdGoGetConfidentialTxData(handle, baseTxHex)
	fmt.Print("CfdGoGetConfidentialTxData end.\n")

	// blind
	fmt.Print("blinding2 start.\n")
	blindHandle, err := cfd.CfdGoInitializeBlindTx(handle)
	if err != nil {
		fmt.Print("CfdGoInitializeBlindTx fail[%s]\n", err.Error())
		return
	}
	defer cfd.CfdGoFreeBlindHandle(handle, blindHandle)

	// 課題：同様の理由から、Signを回数分走らせるのは危険。Blindと同じ形式にすべき。
	// Blindも省メモリ型のを作るべき。
	for i := 0; i < maxTxin; i++ {
		path := fmt.Sprintf("1/%d", (i + 1))
		childKey, err := cfd.CfdGoCreateExtkeyFromParentPath(handle, masterXpriv, path, (int)(cfd.KCfdNetworkTestnet), (int)(cfd.KCfdExtPrivkey))
		if err != nil {
			fmt.Printf("CfdGoCreateExtkeyFromParentPath2 fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}
		blindingKey, _, err := cfd.CfdGoGetPrivkeyFromExtkey(handle, childKey, (int)(cfd.KCfdNetworkTestnet))
		if err != nil {
			fmt.Printf("CfdGoGetPrivkeyFromExtkey fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}
		// unblind
		asset, satoshiAmount, abf, vbf, err := cfd.CfdGoUnblindTxOut(handle, baseTxHex, (uint32)(i), blindingKey)
		if err != nil {
			fmt.Printf("CfdGoUnblindTxOut fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}

		err = cfd.CfdGoAddBlindTxInData(handle, blindHandle, txData.Txid, uint32(i), asset, abf, vbf, satoshiAmount, "", "")
		if err != nil {
			fmt.Printf("CfdGoAddBlindTxInData fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}
	}

	for i := 1; i <= maxTxout; i++ {
		path := fmt.Sprintf("2/%d", i)
		childKey, err := cfd.CfdGoCreateExtkeyFromParentPath(handle, masterXpriv, path, (int)(cfd.KCfdNetworkTestnet), (int)(cfd.KCfdExtPrivkey))
		if err != nil {
			fmt.Printf("CfdGoCreateExtkeyFromParentPath3 fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}
		confidentialKey, err := cfd.CfdGoGetPubkeyFromExtkey(handle, childKey, (int)(cfd.KCfdNetworkTestnet))
		if err != nil {
			fmt.Printf("CfdGoGetPubkeyFromExtkey fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}

		err = cfd.CfdGoAddBlindTxOutData(handle, blindHandle, uint32(i - 1), confidentialKey)
		if err != nil {
			fmt.Printf("CfdGoAddConfidentialTxOut fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}
	}

	outputTxHex, err = cfd.CfdGoFinalizeBlindTx(handle, blindHandle, txHex)
	if err != nil {
		fmt.Printf("CfdGoFinalizeBlindTx fail[%s]\n", err.Error())
		return "", err
	}
	fmt.Print("blinding2 end.\n")
	return outputTxHex, nil
}

func TestCfdCreateLargeBlindTransaction() {
	handle := uintptr(0)
	sequence := (uint32)(cfd.KCfdSequenceLockTimeDisable)
	networkType := (int)(cfd.KCfdNetworkLiquidv1)
	maxTxin := maxnum
	maxTxout := maxnum
	// mnemonic: accuse traffic neglect mechanic sand page cycle tattoo bonus sheriff field top vote outdoor drop
	// master xpriv: tprv8ZgxMBicQKsPe4wfWRxqwNiLTkpLoSvHsKkSP9XXfbV1KVRazP6HNhfDb85caijRwMN8VKFwQu6WeCFjUQGeKk4MTDpHfAahaNHuhUfXhYt
	masterXpriv := "tprv8ZgxMBicQKsPe4wfWRxqwNiLTkpLoSvHsKkSP9XXfbV1KVRazP6HNhfDb85caijRwMN8VKFwQu6WeCFjUQGeKk4MTDpHfAahaNHuhUfXhYt"
	asset := "ef47c42d34de1b06a02212e8061323f50d5f02ceed202f1cb375932aa299f751"

	baseTx, err := CreateBasicBlindTransaction();

	fmt.Print("CfdGoGetConfidentialTxData start.\n")
	txData, err := cfd.CfdGoGetConfidentialTxData(handle, baseTx)
	fmt.Print("CfdGoGetConfidentialTxData end.\n")
	if err != nil {
		return
	}

	txHex, err := cfd.CfdGoInitializeConfidentialTx(handle, uint32(2), uint32(0))
	if err != nil {
		return
	}

	for i := 0; i < maxTxin; i++ {
		txHex, err = cfd.CfdGoAddConfidentialTxIn(
			handle, txHex, txData.Txid, uint32(i), sequence)
		if err != nil {
			fmt.Printf("CfdGoAddConfidentialTxIn fail[%s] idx[%d]\n", err.Error(), i)
			return
		}

		if i % 128 == 0 {
			fmt.Printf(" - txin: %d\n", i)
		}
	}

	descriptor := "pkh(" + masterXpriv + "/0/*)"
	for i := 1; i < maxTxout; i++ {
		bip32DerivationPath := fmt.Sprintf("%d", i)
		descriptorDataList, _, err := cfd.CfdGoParseDescriptor(handle , descriptor, networkType, bip32DerivationPath)
		if err != nil {
			fmt.Print("CfdGoParseDescriptor fail[%s] idx[%d]\n", err.Error(), i)
			return
		}
	
		txHex, err = cfd.CfdGoAddConfidentialTxOut(
			handle, txHex, asset, int64(10000000), "",
			descriptorDataList[0].Address, "", "")
		if err != nil {
			fmt.Printf("CfdGoAddConfidentialTxOut fail[%s] idx[%d]\n", err.Error(), i)
			return
		}

		if i % 128 == 0 {
			fmt.Printf(" - txout: %d\n", i)
		}
	}

	bip32DerivationPath := fmt.Sprintf("%d", maxTxout)
	descriptorDataList, _, err := cfd.CfdGoParseDescriptor(handle , descriptor, networkType, bip32DerivationPath)
	if err != nil {
		fmt.Print("CfdGoParseDescriptor fail[%s] idx[%d]\n", err.Error(), maxTxout)
		return
	}
	txHex, err = cfd.CfdGoAddConfidentialTxOut(
		handle, txHex, asset, int64(9800000), "",
		descriptorDataList[0].Address, "", "")
	if err != nil {
		fmt.Printf("CfdGoAddConfidentialTxOut fail[%s] idx[%d]\n", err.Error(), maxTxout)
		return
	}

	txHex, err = cfd.CfdGoAddConfidentialTxOut(
		handle, txHex, asset, int64(100000), "", "", "", "")
	if err != nil {
		fmt.Printf("CfdGoAddConfidentialTxOut fail[%s] idx[fee]\n", err.Error())
		return
	}

	txHex, err = TestCfdBlindTransaction(txHex, baseTx)
	if err != nil {
		fmt.Printf("TestCfdBlindTransaction fail[%s]\n", err.Error())
		return
	}

	if err != nil {
		errMsg, _ := cfd.CfdGoGetLastErrorMessage(handle)
		fmt.Print("[error message] " + errMsg + "\n")
	}

	fmt.Printf("txHex = %s \n", txHex)
	fmt.Print("TestCfdCreateLargeTransaction test done.\n")
}
