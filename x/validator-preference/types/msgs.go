package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// constants
const (
	TypeMsgSetValidatorSetPreference = "set_validator_set_preference"
)

var _ sdk.Msg = &MsgSetValidatorSetPreference{}

// NewMsgCreateValidatorSetPreference creates a msg to create a validator-set preference.
func NewMsgSetValidatorSetPreference(delegator sdk.AccAddress, preferences []ValidatorPreference) *MsgSetValidatorSetPreference {
	return &MsgSetValidatorSetPreference{
		Delegator:   delegator.String(),
		Preferences: preferences,
	}
}

func (m MsgSetValidatorSetPreference) Route() string { return RouterKey }
func (m MsgSetValidatorSetPreference) Type() string  { return TypeMsgSetValidatorSetPreference }
func (m MsgSetValidatorSetPreference) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Delegator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid delegator address (%s)", err)
	}

	for _, validator := range m.Preferences {
		_, err := sdk.ValAddressFromBech32(validator.ValOperAddress)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid validator address (%s)", err)
		}
	}

	return nil
}

func (m MsgSetValidatorSetPreference) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners takes a create validator-set message and returns the delegator in a byte array.
func (m MsgSetValidatorSetPreference) GetSigners() []sdk.AccAddress {
	delegator, _ := sdk.AccAddressFromBech32(m.Delegator)
	return []sdk.AccAddress{delegator}
}

// constants
const (
	TypeMsgDelegateToValidatorSet = "delegate_to_validator_set"
)

var _ sdk.Msg = &MsgDelegateToValidatorSet{}

// NewMsgMsgStakeToValidatorSet creates a msg to stake to a validator.
func NewMsgMsgStakeToValidatorSet(delegator sdk.AccAddress, coin sdk.Coin) *MsgDelegateToValidatorSet {
	return &MsgDelegateToValidatorSet{
		Delegator: delegator.String(),
		Coin:      coin,
	}
}

func (m MsgDelegateToValidatorSet) Route() string { return RouterKey }
func (m MsgDelegateToValidatorSet) Type() string  { return TypeMsgDelegateToValidatorSet }
func (m MsgDelegateToValidatorSet) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Delegator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	return nil
}

func (m MsgDelegateToValidatorSet) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgDelegateToValidatorSet) GetSigners() []sdk.AccAddress {
	delegator, _ := sdk.AccAddressFromBech32(m.Delegator)
	return []sdk.AccAddress{delegator}
}
