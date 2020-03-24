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
    // TestCfdCreateLargeTransaction()
    // CreateBasicTransaction()
    // TestCfdCreateLargeBlindTransaction()
    // CreateUnblindToBlindTransaction()
    CreateUnblindToBlindTransaction2(false)
    // fmt.Println("end. finishing...")
    // time.Sleep(10 * time.Second) // GC wait
    time.Sleep(1 * time.Second) // GC wait
    fmt.Println("finish")
}

const maxnum = 256
const outMaxnum = 6

func TestCfdCreateLargeTransaction() {
	handle := uintptr(0)
	sequence := (uint32)(cfd.KCfdSequenceLockTimeDisable)
	networkType := (int)(cfd.KCfdNetworkLiquidv1)
	maxTxin := maxnum
	maxTxout := outMaxnum
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
	maxTxout := outMaxnum
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
	maxTxout := outMaxnum
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
		fmt.Printf("CfdGoFinalizeBlindTx2 fail[%s]\n", err.Error())
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
	maxTxout := outMaxnum
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


func CreateUnblindToBlindTransaction() (txHex string, err error) {
	handle := uintptr(0)
	sequence := (uint32)(cfd.KCfdSequenceLockTimeDisable)
	networkType := (int)(cfd.KCfdNetworkLiquidv1)
	maxTxin := maxnum
	maxTxout := outMaxnum
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
	fmt.Printf("tx = %s \n", txHex)
	return txHex, err
}


func CreateUnblindToBlindTransaction2(dump bool) (txHex string, err error) {
	handle := uintptr(0)
	// sequence := (uint32)(cfd.KCfdSequenceLockTimeDisable)
	// networkType := (int)(cfd.KCfdNetworkLiquidv1)
	maxTxin := maxnum
	maxTxout := outMaxnum
	// mnemonic: accuse traffic neglect mechanic sand page cycle tattoo bonus sheriff field top vote outdoor drop
	// master xpriv: tprv8ZgxMBicQKsPe4wfWRxqwNiLTkpLoSvHsKkSP9XXfbV1KVRazP6HNhfDb85caijRwMN8VKFwQu6WeCFjUQGeKk4MTDpHfAahaNHuhUfXhYt
	// masterXpriv := "tprv8ZgxMBicQKsPe4wfWRxqwNiLTkpLoSvHsKkSP9XXfbV1KVRazP6HNhfDb85caijRwMN8VKFwQu6WeCFjUQGeKk4MTDpHfAahaNHuhUfXhYt"
	asset := "5ac9f65c0efcc4775e0baec4ec03abdde22473cd3cf33c0419ca290e0751b225"
	asset2 := "2dcf5a8834645654911964ec3602426fd3b9b4017554d3f9c19403e7fc1411d3"

	txHex = "0200000000fd0001aed7a81c95bf00dcc74b920f6f58696ed69a46eb3846eb9ed175f75d9bf0edcf0000000000ffffffff2cf8b2bd7d26061a2174769399f93653bfc59470070433a913793641a4696cb30000000000ffffffffd579f39c02cfffd442c48d73927d1bfcc25727d73ce3d508b29dd17444c21e110100000000ffffffff99bbb56cf5b1a68a1ec9b55a128c097c3355b8ab24777f5fd6f96e7bee01f2100100000000ffffffffd921931abd9e3b59e283796e27fc5a4d0c14f5ec1ba3da7ab9d71ca1ce1582e20000000000ffffffff7f05c2c4d4bc221bb338fd984094f75fd7ad79d63766c40d9ec179ba758415ec0000000000ffffffffe250a6162e6e19618e48618a76685140cac73986b85ec7aec6e43e443fd24f3e0100000000ffffffffe2dbc93a7d7e6a34cdcce8f378e9163151acf02608529d5d3d4f32c20a1b85f50100000000ffffffff8fce8490a0c352c9bf092360ed1e06edff55e1c3aea6575aaf463581870a69ca0000000000ffffffff280ba7cea2de6fdc4a85c109b6c23c5a8ec12eb2a2db0559afee35ce16ec85b00100000000ffffffff40387f066436e588b3357e0bb3a3bb7e2a899ffe766a5a3d66a536d34a5bf51a0000000000fffffffffbc04c1a2e72c5c2a2291f78e9f2bdfde6a2f7c4faf3d2a197bd166185ebf57f0100000000ffffffffdfe987692a6b1a84169c3956d0cb2003e06be5af7984588c30ece3119b0466330100000000ffffffffde7ba4a68f0df35a0868339e94ab37aaa121c9625212532e72dafb5c81708e830000000000fffffffff6bfaa4e6bec929e8a7a3be1a1a0b9e147e96ade7038f7fc149e978817378ad00000000000ffffffffc0979c57409fac0d3909fb9148475bcc9c2b44e9fd6d91d85275a31d18ac5d060000000000ffffffffb8f010ecadd324048d0c4698812100960b3ef9a0e886a8ed3fff65262ebcf7340000000000ffffffffc89924c6390b87c628b9686e0289370546a4b386acf22ab9120464ef0638c31e0100000000ffffffff5a8e85a8fffb5ce6ff9f4d220acdb47a674e1a5432048a549e6e6934d52eec950100000000fffffffff4471ddccbd63f6088dabca53bea1ebb0c64c0f7f0665211ab7e78c5ceb63eb70000000000ffffffffd99b48f32bda1560cca937b0e30628802eec53139066e01f36b226a250b3abd50000000000ffffffff0de946bc43d428f0c9599577e096accc8c5c51c8130c5d87c90bfe5949e4bc2c0000000000ffffffff2678831e4a85d5364fcb64c1d701080527bd7bad609ea502f5ca276b0fc474b00100000000ffffffffd00b7cbb068d20a6498dd89c2f89766534f6cdcd9ce5c238cf53cce1e8c9cc9b0100000000ffffffff30651a0dc8a8d63a89e0cae1bf0f97df1082c63caeebb6e22adb3d23056de0510100000000ffffffff1d39a41182cb7a3cb6be6acafea21e8bd8d0ef695e6590b0d8c5465d3d1b2eed0100000000ffffffff3c9f17f2fddbaefc53385ea13f6fded9bc16b70e993754655e5800c78b1e4ea40100000000ffffffff97c1bb3811906941cd82b6ef916abeb6322504a03c4a6e4503c17ea4f323b23b0000000000ffffffff639676b00375b92af63aca6cd8c3ff097659507b17dda5de18733f4b8c44bae10000000000ffffffff7cf9b6fe309922bc8f637bb5b7175c32259386cecc749ae3a20b6e44041a873b0100000000fffffffff9f98616e4d16b7a6080c64f72a99df61039441ed0bc2b639079c0c74d1f45c10100000000ffffffff33595fccdd3b7a8d15e8409fbeab091f10db17ff2a7f18da3ee0c3e302ab04920100000000ffffffff2e0261c906976246394804e489bd31021d3dd48e93e25a14813b0dfe2bfb22640100000000ffffffffeb16c0aadb729afca28d701cb147328b3ac22dd19aca061ba3a84ac81da4e7a60000000000fffffffff4a36968b3f67835c6d83b8181a22a234d6d00adaeb8afb5d4e48754248d9e530000000000ffffffff0966f7820f5dbe8571a643283e57bd98bc5c0cd44285ea9734b5f86def46571d0100000000ffffffffca54ced0685a70c7d95a1c92b94ea00a17b4d5b86938a08d34724f86de73017c0000000000ffffffff343a9d19d6decb79605330bcd561f9f96ec54ab1e143d411fc40c0a47034e5680000000000ffffffff3526f6344d051bf356904573caad18c17188773d8e583f60a239f6872e57eccf0100000000ffffffff1ba34c93b9ff6ad722263ee0040924c727bacb6762eed565834cbac9ce111b950000000000ffffffffbee29e1b3dce2fe3aee6f70fd73e700c9a8e4655747749818fb8c1b268fc68b70000000000ffffffff614c0bdb118610e43919af6d8038992aac21cd33be9cbf32b067a77d9507e3ba0100000000ffffffffd042bc3ec77100997343a9aa2da82d090b82a23ff3de34f99f2b5208c722e8450100000000ffffffff3fa3d37a1faa9006ceb1d159691263b0291f14a65b41e06cd03c6133e031d7230000000000ffffffff41530990df981ca25b1bdbd41c285d2778f481474291865a06c38e7a5d35f0500000000000fffffffffde2da00c06c908de5d8db65596d79e7269f1cfbab3602ec3515855c8884d8840100000000fffffffff100794ec29f4eb3d0ecf3810d8c6cb4b96f4fe834bcac8d833438e5a168d02a0000000000ffffffffd62d43fc592551c6ae0afd36f63fdfae45d10087e067897a814d4924e77031b60000000000ffffffff0003cf2babfcd1643a94a35c51f4cf46e92a5d5b3cb3614bbf525cef8d875f260000000000ffffffff2657a1a79bce053aab9c33b2b3b718b80a21bd40c4e6dc439f45373af8c5cd1e0100000000ffffffff754ab5a6a543e8cd3bd5f071be83a03382070eecfd28d2680f70b3c84be0273e0100000000ffffffff047689803896c452e08ea946122ec62daefd04ae9d095939d144abcf46e80d080000000000ffffffffb48f6adb0efff33183cbaf5293114568e5729fb2a1bb7584eeef0aa90a92bfaa0000000000ffffffff550c6e2e79b72830fbf0ffe2d11985d51d183fd9af86376523b75f344e19c81f0000000000ffffffff35f313157d46321a32672b706ac85cec2b78030a3e583812609620668976c59d0000000000ffffffff17f46176eab891b73fb503bd4ff17518c22d4d93cd3b7146116814bc6e23e2f60100000000ffffffffc6aed9bec09210d31ffd27453b90d2997452eaec45e1bcf9ce2a4f1b5c7ef5880000000000ffffffff24ec54a71d3a792faafd413510f539ea091949f1b6681b4823ad4cb320e03af50000000000ffffffffab792c95c3da4075d20291847784ddff7e58dac62e9835d8e040df7e95f0c9240100000000ffffffff5af29ebea69b9cb0814a0dfa8814c3bc446e35dba9392044b1bf2fb56632884a0100000000ffffffffeb36777757f506340e568257dad47627779a08781b6092c73d53f96e063446210000000000ffffffffc5dd73174b247404b2b8bbdc6fe2a1b26729219c87647476b1893fe8368810310000000000ffffffff10e845f3a600e8c05803bc3d181f4d89ba56f15c46e79b202da7ce622561bef40100000000ffffffffcf37297c80bfa3558685ea8fdf692ada90e541b88f9ae4a7278f594d01d7c3360000000000ffffffffdffb755dae53717353a61d29a990c7272a5c8b625718b2c7bfd59412370cd5a40100000000ffffffffa921231d568c8e4cfe19dfbae12ead45274c138fcc8f84a71c5d74dda05fe6db0100000000ffffffffe83b75ea6b89bb04d53e65844beb2af21da8a5b24357a0e1638a1c24938989b60000000000ffffffffd24cb4c3f50a8e7d1ae849914fdb9f853126e234129bf83ba553a4a42b87eb4f0000000000ffffffff12dcfd71250d233039b1f1cd7e6bf8176148f780ddc7c2ed17b708dba8588d340000000000ffffffffe3d4ddf0bd6f350215760d69722ec2d155a250fcdf46f70f7e07d089f4dc14050100000000ffffffff97db77f9d5676450d0168db383c004c5f11d690c1e3747ebfa72a7360f3349c50000000000ffffffffce6f40cbee1bfa4a371c7745ab1ec978b5b99fec2e9112a57fdbde2e5b0631740000000000ffffffff29e55db2e58681618aa1426163a5831cb1b8c7364ef1c7e636e7abc0cd08fc260100000000ffffffff36fd0951020c8e1deac6090771cdc0554905191b471179a08352caa7ebb51f1c0100000000ffffffff563ec7e00bf9be5738dc2d50f1f1c312e1e39f8f8bdf92d4e25b9b84b4bec69b0100000000ffffffff2ad103352c6b3d9080dae2661e6557cd1d70620786b469182fe0068a90752dd10100000000ffffffff5ad370bb7c497ffa1037f7e1684f1205c37ac044c55df3b5223e1432bc725b150100000000ffffffff7d9a04566bd06e932a83899c55a9b0ad7abd4e924c0cc5e830d1acc316e4f63f0100000000fffffffffb43e897561a219e8dcedab1e683ee1099c60406e2249d66b5d6578268bf2c0f0100000000ffffffffbe3a7fba8851fb4282acfb1e4a6796c11a242e97cc3cf7cd0dd121f425b0a1fc0000000000ffffffffc43bf3e278fe67871a0d061aef3e29c28fd8cf4339dd61ced5991c1a0ccbeca30000000000ffffffff6c8c5c861cf1af921eb0199260814a2b9328c3a9bdc232825766cc73dcf1da810000000000ffffffff5a379628849721ba92c827a100c357cc1feba36a493c54201da9e2e262e6b6060000000000ffffffff4d4750eeb852ea398046078701e63eff1e8a20ac141c5ef39033437f16c2414a0000000000ffffffff2c83e4a338ed55df96b2278af6181dff75a6113571413123542b4c9e8c1bbb240100000000ffffffff8e0acc55b4f0eddcf0d7fc51f41fde4baad9db4395c146e4503ec78a489835e20000000000ffffffff75752260f0b03cba4445b59062372b866bbd4b0d9d9f06ed006deeb25f8cf0760000000000ffffffff5989d3ef508f91726f1b709372e7cb96c1856acd9a3203220092a2b48cafcb7d0000000000ffffffffd1072ce7367713fe97950386c93c90c072bfebd659d99a27a065f9d312df4a600000000000ffffffff927cfe7350ab37c3f2022e81905708a72c8c32483bd2a5ae7eb882b5a15c19910000000000ffffffff662a52ed67a4b461885317c31db7d9285ca518e890d674fa8522308f04f752f90100000000ffffffff164fc2737a72fa54b60890f6a5056c23d9aa42da80a23cbd2aa42fe40664142c0000000000ffffffff4ce65427fb8728c810cf59178714d747e9c9cd0dad16d05e8a486b92642db4d40100000000ffffffffdcd68f47bf9ab5c7a0e4a1cc299e2b401fb18eb2ca312c7045d9b29a846184b90000000000ffffffff952135000c7d1c300dd4fed48e6043c9e6b4ce598d6cf7782588de65435170d80100000000ffffffff08c5b1d3f9c2c03deea037a08ed42cebed76f2a4280c4ea4edaff3f53770de080100000000ffffffffad65254b15748d30070130a961e4456f5155ee922bf776a840a75bbe4833bdce0000000000ffffffff9144ec61ea4aa39d007a2ff257ec96535f3c5d29e5423ea280a99a6baeeafe8c0000000000ffffffffff81bfa53dd4d793340bf21224de22c0460aef9d638dc93c2fb14e4b21df815f0000000000ffffffff5b83ea88176cc4f546bceb44b8d56026bbbbbb0fe8be31391725b2b124c3f1d20000000000ffffffff173c01a1fad0b883a580c800952b43668b8665b8ce1c56dea4946d39777da5df0000000000ffffffffebace7b91ff93003721c7f6f596dd8560ecb42409f6a53be7e572755525958b20000000000ffffffffe2044ca2eba03df34275f7ebb61a47c706634cc15fbc6e3d5d0151172a8e4e790100000000ffffffff3b0a9b8e53f866817014caf4ad93ff59e801fc8b3e60a196baee7eba67b6aa5b0100000000fffffffff9f6534dd4c2a4773c7cc4dce7d1f21aedd2e05f3a9a05439d5a390e701d2e3d0100000000ffffffffc6a68336e48518ab2f27260a6d3bf5f5cbefd9ae6af9434034b073be4e89a9730000000000ffffffff84ab5db1c98a44661b46e99620fe036b36abaac2839fcc8517b7f5195b29c8c90000000000ffffffffbc7676461eb6f588e958c202430592cb7d1bcc829ae748f5f06e010b990934a30000000000ffffffffa1b3d91f85acbff118817bc9ff6aed88c5569b06ddb18b09f783de05f448b2530000000000ffffffff898afcc16225196d79e90b7f5a2e5fde6eca3451a512e2c42813c0e85be588eb0000000000ffffffffe05761177ff3a6f5751db4ec9838e08a039169f6d6ea26ae91608a206c36ed630000000000ffffffff2318b2be3aa70b86c47f2a43534cfa77122d6835fc39cebda7ddb4009cabc61b0000000000ffffffffa589f9be43e8805e3f040f62049924b4c8dbb2e8ef9fa2a1b35506eb525695890100000000ffffffff0934c5f7e2f2ce4a9bd7fc1584fa8d3bd18793f08408b967e15cb909d5a612510000000000ffffffff126cd5c183c4c67e7767c8e87d8120eceb0655a982c915d46cc8532e246ecbce0000000000ffffffff2ec89777963a6b645816351aac5e996a518eb6d7b4d80e64c606c48f8f5fea670100000000ffffffffae1efde92f1e68b777c4c328472e295b6dc006f2df702c61b27a2acd8363cda60000000000ffffffffdf310c77ec1f8065b259d77428d8f8bd46aaddc756de96ee8e7ce2356b19e0ef0100000000ffffffff4bde3019da9a02f79e3f1a5446beaebee3b6e1bf4563a71084745b1afe9d9a060100000000ffffffff1705ecc47f28465230d4d3656311a2c13863e42bdca2fcf3b969965c5821883e0000000000ffffffffb04e1c50ac5a9826d4ea2fbb2a7c5b747369bbebd80c2cfae10df73d670cf33b0100000000ffffffff036b58ce797b0cde36f6d6ecf803f8d7c4bf18b6bce52d02eff712bca75d69480100000000ffffffffaf1812c295a972cd8e5aa2f17346685c59cc87b00f96cb7d1bc8a5bccc05b8000100000000ffffffffb276185cd42597f6e3f12461a389dbb05c40f39fe831317c88a6d6a40d0b47440100000000ffffffffde1c0195f87d3a9e1eb70c9b019a587abb791357913e7aa6ff1e9040d86aaa2d0100000000ffffffff217c00a63723b28a7c67ec4760a8e3481bb0e9da82b3a22ad5f8a02a97638efe0100000000ffffffff05f188878022d0a3934949061130d90f06eeb9371a0029a7b082851e68ba030a0100000000ffffffff1f6602098dd142cefaab615a53b8a49a67e337846f321d447ddb02cc9ab523840000000000ffffffffadc86613595fb4dfde069543a4dc91effd6f26d76ac380ba6c3d446339e496c90100000000ffffffff2297dda41b0d9360b99e2c6dc5181e0c56179cfbe0de5181009d99edfafff5af0000000000ffffffff4fb5ad69f97bdbe02fd6c872b55d6552a0fe774aa8788b7cc26a8eebff59acc60100000000ffffffffebbaeb004574f828fed6694ac0bf4945116ce9d584c8bdd6c14226493efed4020100000000ffffffff08550cf94ed40313a5fa8cc97dce9708ff54a66a96006e21fe66bc54efa158820100000000ffffffff85328d19309f53bf3b50f79c61657ad68b03a2f27da60fcae44a80e413264e250000000000ffffffff9dc6874d9448e8accc6890937a10a13071b277b8735a42670eab2a26883de01c0100000000ffffffff64f40d71bcbe167db266347362bd0bdbee130ffb5f21194b3686ea3a47a6d93b0100000000ffffffff5b8554b45844ee21bec534cb8ff0c92077854672ec05c3c0b8cd4e3da0600a7e0000000000ffffffff9438e78de6d01aba8fd981999480124f8f18562d8c9b0658f4f270e7324951470000000000ffffffff5e17432ce51f91d19cc52febba9a88d077cb17f1a905da2679099ee715ba4c230000000000ffffffff4b059384067ebb55f9f46d1426333bcfbeef71989b36deab8cfdc75f1de05f1a0000000000ffffffffa297d4c6defe8abc77b688565f95e1a693a399d87f3aa0a9bf6907647fb36f9d0000000000ffffffff67d774caa3676633dbb196ff7dbff21b59da03bc4ea00ca4d6ef3f82d399f8b30100000000ffffffff1bb7e31f6cfc99b28f030e94998c49ce150c988c2b2c780889ec1531ac95782c0100000000ffffffff58ff23499775d62d443c8c3e7529966bf297edb3be35d7a433a88fd4b73025250000000000ffffffffc06636ea30e5d587ebc88f25ddd5f2a6779e3d7ad090f017a6169b6856b7fb830000000000ffffffff37c627b92dde40216c850d067d7826ee0e3a802f521ee6473847a247ad74835f0000000000ffffffffad233fe495c0a712e084346ffc0920692a9b5c0eb1ed5849f93ce77ba7d992b20100000000ffffffff85fdce15531aa55f4d26f3131259965eb8643b2e402356468da787710717f8210000000000ffffffff2509564cb2731d6107bacec13ff0e5a419121adcfbf82b89ec75fe42e5a750730000000000ffffffff259057a55d2e005d1e850db9600aad2ee0cfc7604bd3d988acef78eb1c2607fa0000000000ffffffff9ec74cb1ca0799ed81f84353be192480412be709bf86fbc5b5c5d723263f4b040000000000ffffffffc5e88af07ec66cdfcd8b62651d7009e381151368d47151b2e8a1bcc43a4b92e10100000000ffffffffce1ad2a1c058abf8cd4dfd85f147bf824d952432f8c77f993f775a44d6c284260000000000ffffffff577f26a1cb0dc313a47ffa0c3ad043b483975f8d37dbab5d3fc450f83f5e8dd00000000000ffffffff71a39af8fa03593cacf21a743460164405e900ab7035b58dce06546ffe80a5810100000000ffffffff51e4e3f812931ef43dc5818e9ed599088eb024eb8dde5b55795ef8b02c2846bc0000000000ffffffff5c3a5cbe66d610efbb0919dd67cd22c2fc4da41276012eca00559c7514f113710000000000ffffffffdd6d99fdc616f47377bd8e8a54933ecefb8fcbd58d7198aad6339219ebfaaae80000000000ffffffffa202e35e32305a25da4dce5ec3a0e15dc7f87a2a1b2bb60747ad4b1f40971eb30000000000fffffffff2976005af43a971fbbf61ba9a7de792c4d1b8892c7853aee05a2a88ff37c0100000000000ffffffffd285ab3e52e495290e1fbdc23353efba4923c9d5de5e6180ab95c66f7a11c68c0100000000ffffffffebca6d7a79e9791eb40377afb505768716a948623d5c360dafb550539dd37cdd0100000000ffffffff6646a477d58e85572d4ff7967052e1a51a8e9369b69b3cd068adc5aba64f3b510100000000ffffffffe131f8e5b98e2aaa3e2fb0ec1c2d49c060fe09d75ff7efbe5aa9e6909c6c69da0000000000ffffffffcd6e35fb4a4ff3bd4c58cafc61217953aac5e2a1f4ab4162a5a11977e406b1f20100000000ffffffff9817a36f0b00bad284c24877fa1274be8120bfad8e8d2c7faf63bfa95632ea0e0000000000ffffffff261cee0f662853faa9c8d45031f2a3cb30514553dc150ae265309d6e2b1894510000000000ffffffff6e24da57d24fb488bd1ab85f7ebd3e987fb728f8762d5729c4c2d08fc1b42bce0000000000ffffffffbc09cb15f486b785bd1d33d25b583b5f7f38655b179cf1551200c0078b0c43e60100000000ffffffffe11b391ffa90fc7f4df24f5695070b3645ae2f02342fc9b6a9b98b994185c0200000000000ffffffffd937300562e5cfe528d0109c0e76089f6c84a665c33220dbd76f997b5516c2db0000000000ffffffff9d5d78e155a0827ef115a0caae68f549adff7900e1a58a85525c550c2dc15fd20000000000ffffffff8443ecfbd35bfb5b9f0107338f5a34837460dee78f14ae00653871ef9e1987e30100000000ffffffffa8541bdc3106519ec2e1bfd854a04679676bc6d59ee95cb91e8e1ac53cdddfa80000000000ffffffff0fd1ebdf711a54d661349ef3dec33f043b8655cd8fcbfc29045707ced55fd10a0000000000ffffffff77e0a35833fcb82ad377c6a57ec130f2d86e9dc8e84722b184d61805c52e8f0e0100000000ffffffff9590b1fe42d78898a21bfbdedfa98145b515e8540b3cf88dc095e1ac883100490000000000ffffffffa15b613ce434772d4ba7e574d7c16a5fbbf693410dedcd04a75295d7381d94460100000000ffffffff7e77b72575f4601de061f0a2cffb8c9dc324a95aea9e3c56f189c79f57bf1c690000000000ffffffffb780d0c091dbb3b6ffae91e9b214991828c4db101455e82a25171ba1973513520100000000ffffffffd0e66541f976f14c331c30063ecbbfd1e36d46f4e3492f21e12eacc5db33a1000100000000ffffffffd0f3c489fe492b2a5c80ff5a490789ca85c495fe52409ca67aba8d6abc8067c80100000000ffffffff307176ccfe6c488e24be8d09f3f462f3bcdf0ca2e928461e17a0c13d53dd40d90000000000ffffffffdd8948e1aa6ad9c1167b4ab1b5d4490a40da2b90b129ee1bc1af8c75aa683d2a0100000000ffffffffbd2da6bf585e3aa5676912280248b8c835934294196f36888f467f773d5ba2e20100000000ffffffff8f8e223bab980368a7499598cf34e9f201a0deba410632aab0e0cd3d8e55e0160000000000ffffffff131b4a3bf1d0e8cbb42ca0ec4b06133de872bbd04115b9dd4bfc164bbd77f8450100000000ffffffffc8fb7fd9428f7e52d2b8d7098d41c72e590b991de788226161ff84ea34e4a8d40000000000ffffffffeb83cb32a1b52a923863b733f75fdc2e0ae9a915fc0832c68690c5a28a4d206b0000000000ffffffff39aff0823de4b6252049e5b7c2b57fcc2a5266a6afbd0fa05d3accec9230529c0100000000ffffffff3a8eeaeb28fcc5f3e3f09e8433de0d82dc40f329dfe6d4231adf62d406c1f32a0100000000ffffffff6ee987963d3cae464aaba992c3521ab204b981f94b885a91d955ff7b00bf1b1c0000000000ffffffff19220721342b970eb5158d1c50e9077d0c9366fcf3fd5650f85de3b3454776420000000000ffffffff8af714eef6bd48a7f96a14091e9b480e6dc7bbf286f14653719ecf8376aa84340100000000ffffffff200c8654eebbfc06ca93a8f05e1160a09cfa8106d5bd113577dbdf1aae0ef2160100000000ffffffffde01e2d594dc52ea21de220c461b18e30013a9f281b18483e97e14f80e5e5fc20100000000ffffffffcdc93753f88b3230250ae0758651bc3f0a71a32788c288a0f450130b943692740000000000ffffffff7713d93e796dedeb737d840a6985b4117fbcb64c93ec74ffc34b2015bb34e4ee0100000000ffffffff0a1def578a54998ad0211399e6ce68b3f5918807484040460706d1401a751cc00100000000ffffffffa4fcef308b83a16bc7abd0bee535538add4cf59dd0e5b2dfdaf3f16797f51f270100000000ffffffff681a7db546df28ffb237ad4138512b2037917f0622e8010457d28ce65518af030000000000fffffffff0d02aa95982ed70f9eaafb9d433e9a5d1b0adcd7887bf36e13e0434ed01d0090100000000ffffffff28f4c16d3dd81cd640c8a572c8108c19503614a5f8cb7d0df9675177bc7feef30000000000ffffffff9fc9349baa04d2472086b53d8d1370c169accefa8b3c89a047a32ad8ecffb9c90100000000ffffffff1f87b0e3b9dcaf5a61c5f7d3efbcc952f06f3b8d0984672125cb2f54079f4e240100000000ffffffff76fbd84683daef78873a81b47454b4b92c762694b6d8911239ba27c073b24d9d0000000000ffffffff401b0f76f226fadb60e685ef49001c051e0428ec4164d961d848ff357a9d556b0100000000ffffffff4f38ad3da9984eb48f2a511d02bcf42466c36480b9ec5ce6b2f36b9555556a680000000000ffffffff5f9b25b20ff82f8b79ae2612e72e948b1b66e27d82503dfcc80821212b884c630100000000fffffffff32612311c816587c777c210d959578c7c9ec5f1abad75ac6beb4f51b51863420000000000ffffffffb4f68affa55427520e856f952453713b102c24d4c6d7b92617ed244ea4b555cb0000000000ffffffff72904588071edc3ca4bb2612841c6d2f61b8e8f673d9ae43f5bd1309f655b79e0000000000ffffffffe90ee6c04cf8546155cf2cb2a2dea7d123f62560e71f5748937fb959df4c69b40100000000fffffffffd1b5f70d410ad0ad694f0f4f926038f4b1679eeb09b114534a3af1c1964b4f40000000000ffffffff25c176be40a1b34986cd9f9cbb35cdfa86de0e4bb1399756a1edc40db0f40d450100000000ffffffff8785c755d7b7daef74ac06634abfb48482db7b8afb4d59fcc2ab71f265536ca50000000000ffffffff490bb2c67164710ddbfb3c50f9cf8534d2fbbcf0818384cc6c10e162a123de710100000000ffffffff4760fa1b9ad0b1c65daeb62c9d8d8d4925ab0b7c8d0ca23f15e82574356504f90000000000ffffffff5fc3b4d6073d68210ae580e4114b403449a7e2cd2a97cc11ffc09078db777a650100000000ffffffffe6df67d7922330de4b440f2cb1fa2eae91e5c19e25c22d17eabf5a7149f674030000000000ffffffff054a366712168fd93b842fa0ebd5469f9be7f28c469c5aa99e4ab1a7644ebefb0100000000fffffffff520f228514ce5986bdda3ca70bbada9b26ae275c3eafc6013c5ce91e1af35a50000000000fffffffff044ea52e108067b61ea49f593f2288196b9fbcf6f92bc1fa0102ee0849622970000000000ffffffff5f128762b8773c4e9a626c2ae4b1022ca42fe0f785407b0088d12168d0e8bec00100000000ffffffffe8feef6645859ff207f8fd285f4725641ad0d5d2b513d78dd63879f56a691c3b0100000000ffffffffbadd7eaffade5b9ea0e2fd692465b5eaa004d19aaf91d29cc5d4f72e09e2f83f0000000000ffffffffa95d643ae2e669410ac002cc54f1535e912bd499ea5e4608e2d65c0548c64a530100000000ffffffff2bbb79f1e6ea084b0c3064c666b278313f2a1c50bd5100eadc0d672cbc8c16e80000000000ffffffff2581ac1f0d4068be6e6b0b92cf3dcf52fb90af942705554fb8b040fd7d4dad730100000000ffffffff6ece82d65bc5626118525b59cb358a6e70b9f7c01daa7456db33173e907499150000000000ffffffff77cf7fad6c712d66d99c09c2f9bcd4e1725782fb754f6485dc19e578b3fce00c0100000000ffffffffbe74bf3d6e774d2493083e8dcb3a2ded58e3fc2a2463abc66a2bea9ad12868fc0100000000ffffffff733a0d59d653470bb473d9adea53e57bd8cf746c85c3a1c96804f603e36e3a6a0100000000ffffffff2acba725dfb1e951f9b712f7c50b6731aae62d6f0cc81754985546e53d86763a0100000000ffffffff9bdf54586964ee71f89506e802e8d0ee76036497248e7e46f9abe16a5babd7a30100000000ffffffff3e180ab4665aee095d41bca8ce4fc901c0636c83435c92b5f5279b47ab8739fb0000000000ffffffff4f94576af42f8c3f1b28f81dc7694b599ec32fe59408016f38b6e5a6d59cca030100000000ffffffff4b5d9a3084cbee0ae13ff6fd185e7cd1c1fd8fb2a3f15d154a2bc8479dd9f0270000000000ffffffffab06fdfd1317845f822e51ce60eb404a38bf5dc533230892e55ff7fb62af4ab00100000000ffffffffb4be462f97312eaca993f3e74104420631afa817084c91bfb79d7d26b7b9dac00100000000ffffffff9903abfb3c5165ced8c8a6e59cca1dda3009f53682e4b90b0ffed318f42601930100000000ffffffff30477d05d8f97d8661ae5b068f142632c5b3f8ebda2643add9469f982272fba20000000000ffffffff938bd4bd0423491fffdbde8f29d9ffd754f04a2c5ea644975e93c16340f24fa60000000000ffffffff6c4460827c83dec3173d143593178017b724062e2362abd7e48dcff4b6f6137f0000000000ffffffff95aceaa0065a297b98c89701aaef19382f69c8598f80577ae7eae39603dd0c8b0000000000ffffffff1429842ab0d11c5e0e1b033054086bb43337ff799f370832246053abb28e4a3f0100000000ffffffffaca8c67964eac1686ee613a74cf86e01804902573b5c063c06f15b583eef147d0000000000ffffffff8c794336409043ea5cd8f6de21192793deee584fa85538d2d75d33a22987ae110100000000ffffffffe0247c7ceb3e1f0883c7bddb3d0f64377f13cff4c8a94916fd8f153d22726c5b0000000000ffffffff299cdf7d1710a4d8a6d0ec448cce3264d106061a44c405c4170eef56934e8ac10100000000ffffffffc739c17078364d5cff619f0403aabe2a768ef5ac3f186b29e5a12180f6d0de550100000000fffffffff21bd5819ec92cf3f4238ffe64e17497f2d5ba988e6c73d3ebe43ffbfb1eff360000000000ffffffff4000a1111b936d4da6e2b091daa6abfcda91ea40e14efb627037804eca080e9e0000000000ffffffff2fb975aea90af30fe463da6b5cf666f4ec9ebc25dc41925ea5dcea9d0f36ba800100000000ffffffff247d7c238d35c52d65e42cded7f6bce8c3edfb142b2c396571ec8cee65a304600300000000fffffffffff63a09f38d965a52903689b60e5321219fa080a7e3c61f936745447ea5f69e0100000000ffffffff0701d31114fce70394c1f9d3547501b4b9d36f420236ec64199154566434885acf2d0100000000000062d4037e3d68c58fa854fab36d479a1a1c2a6ccbf2e2dbfdbe824a7148c48dc005e4d217a9146bde22f6ee71fe7f0871fe0ca641509b51a9a7c9870125b251070e29ca19043cf33ccd7324e2ddab03ecc4ae0b5e77c4fc0e5cf6c95a01000000000000134c037e3d68c58fa854fab36d479a1a1c2a6ccbf2e2dbfdbe824a7148c48dc005e4d21976a914d57f80a95eff6d685e52035e3b9e2aff8db7560f88ac0125b251070e29ca19043cf33ccd7324e2ddab03ecc4ae0b5e77c4fc0e5cf6c95a010000000000269ad0029bde4b5273bd728a9dd8755099ded8f3e9b940f3ba6eedf67070a8a7856c03031976a9140797ac8234831e3e948a14ad42cb40829a45473f88ac01d31114fce70394c1f9d3547501b4b9d36f420236ec64199154566434885acf2d01000000000000c15c029bde4b5273bd728a9dd8755099ded8f3e9b940f3ba6eedf67070a8a7856c030317a914e942a644c891a5847777976589f594f6d0449b18870125b251070e29ca19043cf33ccd7324e2ddab03ecc4ae0b5e77c4fc0e5cf6c95a0100000000000013c402f83876092d7057dff542faff044f6ddf08ae097ae7f31a5ca18090ee8097f6ec1976a914f52b283e344e044003a5cd7974986793dd68808388ac0125b251070e29ca19043cf33ccd7324e2ddab03ecc4ae0b5e77c4fc0e5cf6c95a010000000005f4f5c6039749bd902cfd49894ab58313b01fd66a5e17b8967a2a3d3d1ad3be9a5e21abe31976a914df07adfbbfdf30c50f1105526d9cf959ec28fe9588ac0125b251070e29ca19043cf33ccd7324e2ddab03ecc4ae0b5e77c4fc0e5cf6c95a01000000000000eb3a000000000000"
	// assetAmt := int64(4940 + 2530000 + 5060 + 99939782 + 60218)  // 102540000
	asset2Amt := int64(49500 + 25300)

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
		txid, vout, _, _, err := cfd.CfdGoGetConfidentialTxIn(handle, txHex, uint32(i - 1))
		if err != nil {
			fmt.Printf("CfdGoGetConfidentialTxIn fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}

		//txid := "00000000000000000000000000000000000000000000000000000000" + fmt.Sprintf("%08x", i)
		useAsset := asset
		amt := int64(403600)	// 254 num
		if i == maxTxin {
			useAsset = asset2
			amt = asset2Amt
		}
		if i == maxTxin - 1 {
			amt = int64(25600)
		}
		err = cfd.CfdGoAddBlindTxInData(handle, blindHandle, txid, vout, useAsset, emptyBlinder, emptyBlinder, amt, "", "")
		if err != nil {
			fmt.Printf("CfdGoAddBlindTxInData fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}
	}

	for i := 1; i <= maxTxout; i++ {
		_, _, _, nonce, _, _, _, err := cfd.CfdGoGetConfidentialTxOut(handle, txHex, uint32(i - 1))
		if err != nil {
			fmt.Printf("CfdGoGetConfidentialTxIn fail[%s] idx[%d]\n", err.Error(), i)
			return "", err
		}

		if nonce == "" {
			continue
		}

		err = cfd.CfdGoAddBlindTxOutData(handle, blindHandle, uint32(i - 1), nonce)
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
	if dump {
		fmt.Printf("tx = %s \n", txHex)
	}
	return txHex, err
}
